package database

const InsertURL = "with first_insert as (\n" +
	"    insert into urls(id, correlation_id, original_url,short_url) \n" +
	"   values( nextval('seq_urls'), $2,$3,$4) \n" +
	"   RETURNING id\n)" +
	" insert into user_urls (url_id, user_id) \n" +
	" select id, $1 from first_insert "

const GetURLsByUserID = "select id, user_id, original_url, short_url from urls t1, user_urls t2 where t1.id=t2.url_id and t2.user_id=$1"

const GetOriginalURLByShort = "select original_url from urls where short_url=$1"

const GetOriginalURLByShortForUser = "select original_url from urls t1, user_urls t2 where t1.id=t2.url_id and t2.user_id=$1 and t1.short_url =$2"
