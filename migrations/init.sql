CREATE TABLE "public"."balances" (
    "id" uuid NOT NULL,
    "card_id" uuid NOT NULL,
    "amount" numeric(10,2) DEFAULT '0' NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    CONSTRAINT "balances_card_id" UNIQUE ("card_id"),
    CONSTRAINT "balances_id" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE TABLE "public"."cards" (
    "id" uuid NOT NULL,
    "external_id" integer NOT NULL,
    "name" character varying(80) NOT NULL,
    "enabled" boolean NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    CONSTRAINT "cards_external_id" UNIQUE ("external_id"),
    CONSTRAINT "cards_pk" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE TABLE "public"."deposits" (
    "id" uuid NOT NULL,
    "external_id" uuid NOT NULL,
    "card_id" uuid NOT NULL,
    "amount" numeric(10,2) NOT NULL,
    "paid" boolean DEFAULT true NOT NULL,
    "cancelled" boolean DEFAULT false NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    CONSTRAINT "deposits_external_id" UNIQUE ("external_id"),
    CONSTRAINT "deposits_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TABLE "public"."transfers" (
    "id" uuid NOT NULL,
    "external_id" uuid NOT NULL,
    "src_card_id" uuid NOT NULL,
    "dst_card_id" uuid NOT NULL,
    "amount" numeric(10,2) NOT NULL,
    "cancelled" boolean DEFAULT false NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    CONSTRAINT "transfers_external_id" UNIQUE ("external_id"),
    CONSTRAINT "transfers_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TYPE reference AS ENUM ('deposit', 'transfer', 'purchase');

CREATE TABLE "public"."ledgers" (
    "id" uuid NOT NULL,
    "card_id" uuid NOT NULL,
    "reference" reference NOT NULL,
    "reference_id" uuid NOT NULL,
    "amount" numeric(10,2) NOT NULL,
    "created_at" timestamptz NOT NULL,
    CONSTRAINT "ledger_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "ledger_card_id" ON "public"."ledgers" USING btree ("card_id");


ALTER TABLE ONLY "public"."deposits" ADD CONSTRAINT "deposits_card_id_fkey" FOREIGN KEY (card_id) REFERENCES cards(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."transfers" ADD CONSTRAINT "transfers_src_card_id_fkey" FOREIGN KEY (src_card_id) REFERENCES cards(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."transfers" ADD CONSTRAINT "transfers_dst_card_id_fkey" FOREIGN KEY (dst_card_id) REFERENCES cards(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."ledgers" ADD CONSTRAINT "ledger_card_id_fkey" FOREIGN KEY (card_id) REFERENCES cards(id) NOT DEFERRABLE;
