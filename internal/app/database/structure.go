package database

const CreateDatabaseStructure = "create table if not exists  user_urls (id numeric primary key, user_id varchar, short_url varchar, original_url varchar);" +
	"create sequence if not exists seq_user_urls increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by user_urls.id;" +
	"create index if not exists user_url_user_id_idx on user_urls (user_id);" +
	"create index if not exists user_url_short_url_idx on user_urls (short_url);" +
	"create unique index if not exists user_url_original_url_udx on user_urls (user_id, original_url);" +
	"create table if not exists  user_urls2 (id numeric primary key, user_id varchar, correlation_id varchar, short_url varchar, original_url varchar);" +
	"create sequence if not exists seq_user_urls2 increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by user_urls.id;" +
	"create index if not exists user_url2_user_id_idx on user_urls (user_id);" +
	"create index if not exists user_url2_short_url_idx on user_urls (short_url);" +
	"create unique index if not exists user_url2_udx on user_urls2 (user_id, correlation_id, original_url);"

const ClearDatabaseStructure = "drop table if exists user_urls cascade;\n" +
	"drop table if exists user_urls2 cascade;"
