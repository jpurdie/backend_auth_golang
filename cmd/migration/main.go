package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbInsert := `
	INSERT INTO public.roles VALUES (500, 500, 'OWNER', TRUE);
	INSERT INTO public.roles VALUES (400, 400, 'SUPERUSER', FALSE);
	INSERT INTO public.roles VALUES (300, 300, 'ADMIN', TRUE);
	INSERT INTO public.roles VALUES (200, 200, 'SUPERVISOR', FALSE);
	INSERT INTO public.roles VALUES (100, 100, 'USER', TRUE);

	INSERT INTO public.project_statuses (name, uuid) VALUES ('Not Started', uuid_generate_v4());
	INSERT INTO public.project_statuses (name, uuid) VALUES ('In Progress', uuid_generate_v4());
	INSERT INTO public.project_statuses (name, uuid) VALUES ('Done', uuid_generate_v4());
	INSERT INTO public.project_statuses (name, uuid) VALUES ('Duplicate', uuid_generate_v4());
	INSERT INTO public.project_statuses (name, uuid) VALUES ('On Hold', uuid_generate_v4());
	INSERT INTO public.project_statuses (name, uuid) VALUES ('Stopped', uuid_generate_v4());
	INSERT INTO public.project_statuses (name, uuid) VALUES ('Deferred', uuid_generate_v4());

	INSERT INTO public.project_types (name, uuid) VALUES ('Operational', uuid_generate_v4());
	INSERT INTO public.project_types (name, uuid) VALUES ('Service', uuid_generate_v4());
	INSERT INTO public.project_types (name, uuid) VALUES ('Project', uuid_generate_v4());

	INSERT INTO public.project_complexities (name, weight, uuid) VALUES ('Average', 100, uuid_generate_v4());
	INSERT INTO public.project_complexities (name, weight, uuid) VALUES ('Hard', 200, uuid_generate_v4());
	INSERT INTO public.project_complexities (name, weight, uuid) VALUES ('Very Difficult', 300, uuid_generate_v4());
	INSERT INTO public.project_complexities (name, weight, uuid) VALUES ('Extreme but known', 400, uuid_generate_v4());
	INSERT INTO public.project_complexities (name, weight, uuid) VALUES ('Extreme and unknown', 500, uuid_generate_v4());

	INSERT INTO public.project_sizes (name, weight, uuid) VALUES ('Extra Small', 100, uuid_generate_v4());
	INSERT INTO public.project_sizes (name, weight, uuid) VALUES ('Small', 200, uuid_generate_v4());
	INSERT INTO public.project_sizes (name, weight, uuid) VALUES ('Medium', 300, uuid_generate_v4());
	INSERT INTO public.project_sizes (name, weight, uuid) VALUES ('Large', 400, uuid_generate_v4());
	INSERT INTO public.project_sizes (name, weight, uuid) VALUES ('Extra Large', 500, uuid_generate_v4());
	INSERT INTO public.project_sizes (name, weight, uuid) VALUES ('Epic', 600, uuid_generate_v4());

	CREATE INDEX orgs_uuid on organizations(uuid);
	CREATE INDEX projects_name on projects(name);
	CREATE INDEX projects_uuid on projects(uuid);
	CREATE INDEX profiles_uuid on profiles(uuid);
	CREATE INDEX roles_uuid on roles(name);
	CREATE INDEX users_uuid on users(uuid);
	CREATE INDEX users_externalID on users(external_id);
	CREATE INDEX users_email on users(email);
	CREATE INDEX project_statuses_uuid on project_statuses(uuid);
	`
	dbInsert = ``
	var psn = os.Getenv("DATABASE_URL")
	queries := strings.Split(dbInsert, ";")

	u, err := pg.ParseURL(psn)
	checkErr(err)
	db := pg.Connect(u)
	_, err = db.Exec("SELECT 1")
	checkErr(err)
	//createSchema(db, &authapi.Organization{})
	//createSchema(db, &authapi.Role{})
	//createSchema(db, &authapi.User{})
	//createSchema(db, &authapi.Profile{})
	//createSchema(db, &authapi.Invitation{})
	//createSchema(db, &model.ProjectStatus{})
	//createSchema(db, &model.ProjectType{})
	//createSchema(db, &model.ProjectComplexity{})
	//createSchema(db, &model.ProjectSize{})
	//createSchema(db, &model.Project{})
	//createSchema(db, &model.StrategicAlignment{})
	//createSchema(db, &model.ProjectStrategicAlignment{})
	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		checkErr(err)
	}

	//sec := secure.New(1, nil)

	//userInsert := `INSERT INTO public.users (id, created_at, updated_at, first_name, last_name, password, email, active, role_id, company_id, location_id) VALUES (1, now(),now(), 'Admin', 'admin', '%s', 'johndoe@mail.com', true, 100, 1, 1);`
	//_, err = db.Exec(fmt.Sprintf(userInsert, sec.Hash("admin")))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		}))
	}
}
