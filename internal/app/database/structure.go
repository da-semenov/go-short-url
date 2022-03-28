package database

const urls = "create table if not exists  urls (id numeric primary key, correlation_id varchar, original_url varchar, short_url varchar);\n" +
	"create sequence if not exists seq_urls increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by urls.id;\n" +
	"create index if not exists urls_short_url_idx on urls (short_url);\n" +
	"create unique index if not exists urls_udx on urls (original_url);\n"

const userURLs = "create table if not exists  user_urls (user_id varchar, url_id numeric);\n" +
	"create unique index if not exists user_url_idx1 on user_urls (user_id, url_id);\n"

const CreateDatabaseStructure = urls + userURLs

const ClearDatabaseStructure = "drop table if exists user_urls cascade;\n" +
	"drop table if exists urls cascade;\n"
