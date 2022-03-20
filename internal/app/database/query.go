package database

const GetURLsByUserID = "select id, user_id, short_url, original_url from user_urls where user_id=$1"
const InsertUserURL = "insert into user_urls (id, user_id, short_url, original_url) \n\t" +
	"values ( nextval('seq_user_urls'), $1, $2,$3)\n" +
	"on conflict (user_id, original_url) do update set short_url = excluded.short_url"
const InsertUserURL2 = "insert into user_urls2 (id, user_id,correlation_id, original_url, short_url) \n\t" +
	"values ( nextval('seq_user_urls'), $1, $2,$3, $4)\n" +
	"on conflict (user_id, correlation_id, original_url) do update set short_url = excluded.short_url"

const AllUserURLsWithCorrelationIDByUserID = "select correlation_id, original_url, short_url from user_urls2 where correlation_id is not null and user_id=$1"
