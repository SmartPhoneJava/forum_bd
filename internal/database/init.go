package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/SmartPhoneJava/forum_bd/internal/config"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {

	// for local launch
	if os.Getenv(CDB.URL) == "" {
		os.Setenv(CDB.URL, "user=rolepade password=github.com/SmartPhoneJava/forum_bd dbname=escabase sslmode=disable")
	}

	var database *sql.DB
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		debug("database/Init cant open:" + err.Error())
		return
	}

	db = &DataBase{
		Db: database,
	}
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = db.Db.Ping(); err != nil {
		debug("database/Init cant access:" + err.Error())
		return
	}
	debug("database/Init open")
	if err = db.createTables(); err != nil {
		return
	}
	debug("database/Init done")

	return
}

// CreateTables creates table
func (db *DataBase) createTables() error {
	sqlStatement := dropTables() +
		userCreateTable() + forumCreateTable() +
		userInForumCreateTable() +
		threadCreateTable() + postCreateTable() +
		voteCreateTable() + statusCreateTable()
	fmt.Println(sqlStatement)
	_, err := db.Db.Exec(sqlStatement)

	if err != nil {
		debug("database/init - fail:" + err.Error())
	}
	return err
}

// drop tables and indexes
func dropTables() string {
	return `
    DROP TABLE IF EXISTS Vote CASCADE;
    DROP TABLE IF EXISTS Post50 CASCADE;
    DROP TABLE IF EXISTS Post100 CASCADE;
    DROP TABLE IF EXISTS Post CASCADE;
    DROP TABLE IF EXISTS Thread CASCADE;
    DROP TABLE IF EXISTS Thread_forum_below_50 CASCADE;
    DROP TABLE IF EXISTS Thread_forum_below_100 CASCADE;
    DROP TABLE IF EXISTS Thread_forum_below_150 CASCADE;
    DROP TABLE IF EXISTS Thread_forum_below_200 CASCADE;
    DROP TABLE IF EXISTS Thread_forum_other CASCADE;
    DROP TABLE IF EXISTS UserInForum CASCADE;
    DROP TABLE IF EXISTS Forum CASCADE;
    DROP TABLE IF EXISTS UserForum CASCADE;
    DROP TABLE IF EXISTS Status CASCADE;

    `
}

func userCreateTable() string {
	return `
    CREATE Table UserForum (
        id SERIAL PRIMARY KEY,
        nickname text NOT NULL UNIQUE collate "C",
        fullname text NOT NULL,
        email text UNIQUE NOT NULL,
        about text 
    );

    ALTER TABLE UserForum OWNER TO rolepade;

    CREATE INDEX userforum_lower_nickname ON 
        UserForum (lower(nickname));
    
    `
}

func userInForumCreateTable() string {
	return `
    -- there are lower versions
    CREATE Table UserInForum (
        id SERIAL PRIMARY KEY,
        nickname text NOT NULL, --nickname
        forum text NOT NULL
    );

    ALTER TABLE UserInForum OWNER TO rolepade;

    
    CREATE OR REPLACE FUNCTION trigger_userinforum_before_function () 
    RETURNS trigger AS $$ 
    DECLARE
        found int;
    BEGIN 
        select count(1) from UserInForum where nickname = NEW.nickname and forum = NEW.forum into found;
        if (found <> 1) then
            return NEW;
        ELSE
            return NULL;
        END IF;
    END; 
    $$ LANGUAGE  plpgsql;

    ALTER FUNCTION trigger_userinforum_before_function() OWNER TO rolepade;
    
    CREATE TRIGGER userinforum_before_trigger 
    BEFORE INSERT ON UserInForum FOR EACH ROW 
    EXECUTE PROCEDURE trigger_userinforum_before_function ();
    

    CREATE INDEX user_in_forum ON 
        UserInForum (nickname, forum);
    ALTER INDEX user_in_forum OWNER TO rolepade;
    
    `
}

func forumCreateTable() string {
	return `
    CREATE Table Forum (
        id SERIAL PRIMARY KEY,
        posts int default 0,
        slug text not null UNIQUE,
        threads int default 0,
        title text not null,
        user_nickname text not null
    );

    ALTER TABLE Forum OWNER TO rolepade;

    ALTER TABLE Forum
        ADD CONSTRAINT forum_user
        FOREIGN KEY (user_nickname)
        REFERENCES UserForum(nickname)
            ON DELETE CASCADE;
    
    CREATE INDEX forum_lower_user_nickname ON 
        Forum using btree (lower(user_nickname));

    ALTER INDEX forum_lower_user_nickname OWNER TO rolepade;

    CREATE UNIQUE INDEX forum_lower_slug ON 
        Forum (lower(slug));

    ALTER INDEX forum_lower_slug OWNER TO rolepade;
    `
}

func threadCreateTable() string {
	return `
    CREATE Table Thread (
        id SERIAL PRIMARY KEY,
        author text not null,
        forum text not null,
        forum_id int not null,
        message text not null,
        created    TIMESTAMPTZ,
        title text not null,
        votes int default 0,
        slug text default null
    );

    ALTER TABLE Thread OWNER TO rolepade;

    ALTER TABLE Thread
    ADD CONSTRAINT thread_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    ALTER TABLE Thread
        ADD CONSTRAINT thread_forum_id
        FOREIGN KEY (forum_id)
        REFERENCES Forum(id)
            ON DELETE CASCADE;

    ALTER TABLE Thread
    ADD CONSTRAINT thread_forum
    FOREIGN KEY (forum)
    REFERENCES Forum(slug)
        ON DELETE CASCADE;
/*
        CREATE Table Thread_forum_below_50 (
            CHECK(forum_id < 50)
        ) INHERITS (Thread);

        ALTER TABLE Thread_forum_below_50 OWNER TO rolepade;

        CREATE Table Thread_forum_below_100 (
            CHECK(forum_id >= 50 and forum_id < 100)
        ) INHERITS (Thread);

        ALTER TABLE Thread_forum_below_100 OWNER TO rolepade;

        CREATE Table Thread_forum_below_150 (
            CHECK(forum_id >= 100 and forum_id < 150)
        ) INHERITS (Thread);

        ALTER TABLE Thread_forum_below_150 OWNER TO rolepade;

        CREATE Table Thread_forum_below_200 (
            CHECK(forum_id >= 150 and forum_id < 200)
        ) INHERITS (Thread);

        ALTER TABLE Thread_forum_below_200 OWNER TO rolepade;

        CREATE Table Thread_forum_other(
            CHECK(forum_id >= 200)
        ) INHERITS (Thread);

        ALTER TABLE Thread_forum_other OWNER TO rolepade;
        
    CREATE INDEX thread_lower_author_50 ON 
        Thread_forum_below_50 using btree (lower(author));

    ALTER INDEX thread_lower_author_50 OWNER TO rolepade;

    CREATE INDEX thread_lower_slug_50 ON 
        Thread_forum_below_50 using btree (lower(slug));

    ALTER INDEX thread_lower_slug_50 OWNER TO rolepade;

    CREATE INDEX thread_lower_forum_id_50 ON 
        Thread_forum_below_50 using btree (forum_id);

    ALTER INDEX thread_lower_forum_id_50 OWNER TO rolepade;

    CREATE INDEX thread_lower_author_100 ON 
        Thread_forum_below_100 using btree (lower(author));

    ALTER INDEX thread_lower_author_100 OWNER TO rolepade;

    CREATE INDEX thread_lower_slug_100 ON 
        Thread_forum_below_100 using btree (lower(slug));

    ALTER INDEX thread_lower_slug_100 OWNER TO rolepade;

    CREATE INDEX thread_lower_forum_id_100 ON 
        Thread_forum_below_100 using btree (forum_id);

    ALTER INDEX thread_lower_forum_id_100 OWNER TO rolepade;
    
    CREATE INDEX thread_lower_author_150 ON 
        Thread_forum_below_150 using btree (lower(author));

    ALTER INDEX thread_lower_author_150 OWNER TO rolepade;

    CREATE INDEX thread_lower_slug_150 ON 
        Thread_forum_below_150 using btree (lower(slug));

    ALTER INDEX thread_lower_slug_150 OWNER TO rolepade;

    CREATE INDEX thread_lower_forum_id_150 ON 
        Thread_forum_below_150 using btree (forum_id);

    ALTER INDEX thread_lower_forum_id_150 OWNER TO rolepade;

    CREATE INDEX thread_lower_author_200 ON 
        Thread_forum_below_200 using btree (lower(author));

    ALTER INDEX thread_lower_author_200 OWNER TO rolepade;

    CREATE INDEX thread_lower_slug_200 ON 
        Thread_forum_below_200 using btree (lower(slug));

    ALTER INDEX thread_lower_slug_200 OWNER TO rolepade;

    CREATE INDEX thread_lower_forum_id_200 ON 
        Thread_forum_below_200 using btree (forum_id);

    ALTER INDEX thread_lower_forum_id_200 OWNER TO rolepade;

    CREATE INDEX thread_lower_author_other ON 
        Thread_forum_other using btree (lower(author));

    ALTER INDEX thread_lower_author_other OWNER TO rolepade;

    CREATE INDEX thread_lower_slug_other ON 
        Thread_forum_other using btree (lower(slug));

    ALTER INDEX thread_lower_slug_other OWNER TO rolepade;

    CREATE INDEX thread_lower_forum_id_other ON 
        Thread_forum_other using btree (forum_id);

    ALTER INDEX thread_lower_forum_id_other OWNER TO rolepade;
    */
        CREATE OR REPLACE FUNCTION trigger_thread_after_function () 
        RETURNS trigger AS $$ 
        DECLARE
        BEGIN 
            UPDATE Forum set threads=threads+1 where slug=NEW.forum;
            UPDATE Status set thread = thread + 1;
            --INSERT into UserInForum(nickname, forum) values (lower(NEW.author), lower(NEW.forum));
            return NEW;
        END; 
        $$ LANGUAGE  plpgsql;

        ALTER FUNCTION trigger_thread_after_function() OWNER TO rolepade;
        /*
        CREATE TRIGGER thread_after_trigger_50 
        AFTER INSERT ON Thread_forum_below_50 FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_after_function ();

        --DROP TRIGGER IF EXISTS thread_after_trigger_100 on Thread_forum_below_100; 
        CREATE TRIGGER thread_after_trigger_100 
        AFTER INSERT ON Thread_forum_below_100 FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_after_function ();

        --DROP TRIGGER IF EXISTS thread_after_trigger_150 on Thread_forum_below_150; 
        CREATE TRIGGER thread_after_trigger_150 
        AFTER INSERT ON Thread_forum_below_150 FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_after_function ();

        --DROP TRIGGER IF EXISTS thread_after_trigger_200 on Thread_forum_below_200; 
        CREATE TRIGGER thread_after_trigger_200 
        AFTER INSERT ON Thread_forum_below_200 FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_after_function ();

        --DROP TRIGGER IF EXISTS thread_after_trigger_other on Thread_forum_other; 
        CREATE TRIGGER thread_after_trigger_other 
        AFTER INSERT ON Thread_forum_other FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_after_function ();
        */

        CREATE TRIGGER thread_after_trigger
        AFTER INSERT ON Thread FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_after_function ();

        CREATE OR REPLACE FUNCTION trigger_thread_before_function () 
        RETURNS trigger AS $$ 
        DECLARE
            --fauthor text = '';
            --fcreated TIMESTAMPTZ;
            --fmessage text;
            --ftitle text;
            forum_id int;
            forum_slug text;
        BEGIN 
        /*
            if NEW.slug <> '' then
                SELECT author, created, message, title 
                    from Thread where lower(slug) like lower(NEW.slug) limit 1 into
                    fauthor, fcreated, fmessage, ftitle;
            END IF;

            if fauthor <> '' then
                NEW.Author = fauthor;
                NEW.Created = fcreated;
                NEW.Author = fauthor;
                --RAISE EXCEPTION 'error thread exist';
            END IF;
            */

            select id, slug from Forum where lower(slug) like lower(NEW.forum) into forum_id, forum_slug;
            
            NEW.forum_id = forum_id;
            NEW.forum = forum_slug;
            /*
            if (forum_id < 50) then
                insert into Thread_forum_below_50 VALUES (NEW.*);
            elseif (forum_id >= 50 and forum_id < 100) then
                insert into Thread_forum_below_100 VALUES (NEW.*);
            elseif (forum_id >= 100 and forum_id < 150) then
                insert into Thread_forum_below_150 VALUES (NEW.*);
            elseif (forum_id >= 150 and forum_id < 200) then
                insert into Thread_forum_below_200 VALUES (NEW.*);
            else
                insert into Thread_forum_other VALUES (NEW.*);
            end if;
            return NULL;
            */
            RETURN NEW;
        END; 
        $$ LANGUAGE  plpgsql;

        ALTER FUNCTION trigger_thread_before_function() OWNER TO rolepade;
        
        --DROP TRIGGER IF EXISTS thread_before_trigger on Thread; 
        CREATE TRIGGER thread_before_trigger 
        BEFORE INSERT ON Thread FOR EACH ROW 
        EXECUTE PROCEDURE trigger_thread_before_function ();
        

    `
}

func postCreateTable() string {
	return `
    CREATE Table Post (
        id SERIAL PRIMARY KEY,
        author text not null,
        forum text,
        message text not null,
        created    TIMESTAMPTZ,
        isEdited boolean default false,
        thread int,
        parent int,
        path text not null default '0'
    );

    ALTER Table Post OWNER TO rolepade;

    ALTER TABLE Post
    ADD CONSTRAINT post_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    --ALTER TABLE Post
    --ADD CONSTRAINT post_forum
    --FOREIGN KEY (forum)
    --REFERENCES Forum(slug)
    --    ON DELETE CASCADE;

    CREATE INDEX post_forum ON 
        Post using btree (forum);
    
    --ALTER TABLE Post
    --ADD CONSTRAINT post_thread
    --FOREIGN KEY (thread)
    --REFERENCES Thread(id)
    --    ON DELETE CASCADE;

    CREATE INDEX post_thread ON 
            Post using btree (thread);

    CREATE INDEX post_lower_author ON 
            Post using btree (lower(author));

/*
     CREATE Table Post50 (
            CHECK(thread < 50)
        ) INHERITS (Post);
    CREATE Table Post100 (
            CHECK(thread >= 50)
        ) INHERITS (Post);

        ALTER TABLE Post50
        ADD CONSTRAINT post_user
        FOREIGN KEY (author)
        REFERENCES UserForum(nickname)
            ON DELETE CASCADE;
    
        ALTER TABLE Post50
        ADD CONSTRAINT post_forum
        FOREIGN KEY (forum)
        REFERENCES Forum(slug)
            ON DELETE CASCADE;
    
        ALTER TABLE Post100
        ADD CONSTRAINT post_user
        FOREIGN KEY (author)
        REFERENCES UserForum(nickname)
            ON DELETE CASCADE;
    
        ALTER TABLE Post100
        ADD CONSTRAINT post_forum
        FOREIGN KEY (forum)
        REFERENCES Forum(slug)
            ON DELETE CASCADE;

    CREATE INDEX post50_thread ON 
            Post50 using btree (thread);

    CREATE INDEX post50_lower_author ON 
            Post50 using btree (lower(author));

    CREATE INDEX post100_thread ON 
            Post100 using btree (thread);

    CREATE INDEX post100_lower_author ON 
            Post100 using btree (lower(author));
    */
    
    /*
    CREATE OR REPLACE FUNCTION  update_post_path() 
    RETURNS trigger AS $$ 
    DECLARE
        check_user int;
        check_parent int;
        parent_path text;
    BEGIN 
        --select 1 from UserForum where lower(nickname) like lower(NEW.author) limit 1 into check_user; 
        
        --if check_user = null then
        --    RAISE EXCEPTION 'error with check_user';
        --END IF;

        if NEW.parent <> 0 then
            select count(1) from Post as P where id = NEW.parent and thread = NEW.thread limit 1 into check_parent;

            if check_parent <> 1 then
                RAISE EXCEPTION 'Parent post was created in another thread';
            END IF;
        END IF;

        if NEW.parent = 0 then
            parent_path = '0';
        ELSE
            select path from Post where id = NEW.parent into parent_path;
        END IF;
        NEW.path = parent_path || '.' || (NEW.id::text);
        --INSERT into UserInForum(nickname, forum) values (lower(NEW.author), lower(NEW.forum));
        return NEW;
    END; 
    $$ LANGUAGE  plpgsql;
    */

    CREATE OR REPLACE FUNCTION  trigger_post_before_function() 
    RETURNS trigger AS $$ 
    DECLARE
        check_user int;
        check_parent int;
        parent_path text;
    BEGIN 
        --select 1 from UserForum where lower(nickname) like lower(NEW.author) limit 1 into check_user; 
        
        --if check_user = null then
        --    RAISE EXCEPTION 'error with check_user';
        --END IF;

        if NEW.parent <> 0 then
            select count(1) from Post as P where id = NEW.parent and thread = NEW.thread limit 1 into check_parent;

            if check_parent <> 1 then
                RAISE EXCEPTION 'Parent post was created in another thread';
            END IF;
        END IF;

        if NEW.parent = 0 then
            parent_path = '0';
        ELSE
            select path from Post where id = NEW.parent into parent_path;
        END IF;
        NEW.path = parent_path || '.' || (NEW.id::text);
        --INSERT into UserInForum(nickname, forum) values (lower(NEW.author), lower(NEW.forum));
        return NEW;
    END; 
    $$ LANGUAGE  plpgsql;

    ALTER FUNCTION trigger_post_before_function() OWNER TO rolepade;
    
    --DROP TRIGGER IF EXISTS post_trigger_before on Post; 
    CREATE TRIGGER post_trigger_before
    BEFORE INSERT ON Post FOR EACH ROW 
    EXECUTE PROCEDURE trigger_post_before_function ();
/*
    CREATE TRIGGER post50_trigger_before
    BEFORE INSERT ON Post50 FOR EACH ROW 
    EXECUTE PROCEDURE trigger_post_before_function ();

    CREATE TRIGGER post100_trigger_before
    BEFORE INSERT ON Post100 FOR EACH ROW 
    EXECUTE PROCEDURE trigger_post_before_function ();
*/
    `
}

func voteCreateTable() string {
	return `
    CREATE Table Vote (
        id SERIAL PRIMARY KEY,
        author text not null,
        thread int not null,
        isEdited bool default false,
        voice int default 0,
        old_voice int default 0
    );

    ALTER Table Vote OWNER TO rolepade;

    ALTER TABLE Vote
    ADD CONSTRAINT vote_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    CREATE INDEX vote_thread ON 
        Vote using btree (thread);

    --ALTER TABLE Vote
    --ADD CONSTRAINT vote_thread
    --FOREIGN KEY (thread)
    --REFERENCES Thread(id)
    --    ON DELETE CASCADE;

    CREATE UNIQUE INDEX IF NOT EXISTS vote_thread_author ON 
        Vote (thread, author);

    ALTER Index vote_thread_author OWNER TO rolepade;

    CREATE OR REPLACE FUNCTION trigger_vote_function () 
        RETURNS trigger AS $$ 
        DECLARE
        BEGIN 

        UPDATE Thread set votes = votes + NEW.voice - NEW.old_voice
            where id = NEW.thread;
        return NEW;
        END; 
        $$ LANGUAGE  plpgsql;

        ALTER FUNCTION trigger_vote_function() OWNER TO rolepade;
        
        --DROP TRIGGER IF EXISTS vote_trigger on Vote; 
        CREATE TRIGGER vote_trigger 
        AFTER UPDATE OR INSERT ON Vote FOR EACH ROW 
        EXECUTE PROCEDURE trigger_vote_function ();
    `
}

func statusCreateTable() string {
	return `
    CREATE Table Status (
        Forum  int default 0,
        Post   int default 0,
        Thread int default 0,
        Users   int default 0
    );

    ALTER Table Status OWNER TO rolepade;

    INSERT INTO Status(Post) VALUES (0) 
						 
    `
}
