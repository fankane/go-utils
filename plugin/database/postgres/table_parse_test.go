package postgres

import (
	"fmt"
	"github.com/fankane/go-utils/str"
	"testing"
)

func TestTableColumnsFromDDL(t *testing.T) {
	ddlT := `  
CREATE TABLE "public"."tt1" (
  "id" int4 NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "score" numeric(4,3),
  "geom" polygon,
  "t2" text COLLATE "pg_catalog"."default",
  "t3" timestamp(6),
  "t4" bool,
  "t5" float4,
  "t6" bit(32),
  "t7" numeric(255,4),
  "t_array_1" text[] COLLATE "pg_catalog"."default",
  "t_array_2" text[][] COLLATE "pg_catalog"."default",
  "di_id" int4,
  CONSTRAINT "tt1_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "rID" FOREIGN KEY ("di_id") REFERENCES "public"."database_infos" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "unih" UNIQUE ("t2"),
  CONSTRAINT "uniUni" UNIQUE ("score", "name"),
  CONSTRAINT "c1" CHECK (8::double precision < t5)
)
;

ALTER TABLE "public"."tt1" 
  OWNER TO "postgres";

CREATE INDEX "idxna1" ON "public"."tt1" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

COMMENT ON COLUMN "public"."tt1"."id" IS 'auto increase ';

COMMENT ON COLUMN "public"."tt1"."t2" IS 'tttt';

COMMENT ON COLUMN "public"."tt1"."t4" IS 'ffff';   
    `

	res, err := TableColumnsFromDDL(ddlT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str.ToJSON(res))
}
