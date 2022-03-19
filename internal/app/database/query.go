package database

const GetURLsByUserID = "select id, user_id, short_url, original_url from user_urls where user_id=$1"
const InsertUserURL = "insert into user_urls (id, user_id, short_url, original_url) \n\t" +
	"values ( nextval('seq_user_urls'), $1, $2,$3)\n" +
	"on conflict (user_id, original_url) do update set short_url = excluded.short_url"
