package database

const InsertURL = "with first_insert as (insert into urls(id, correlation_id, original_url,short_url) values(nextval('seq_urls'), $2,$3,$4) RETURNING id\n)" +
	"insert into user_urls (url_id, user_id) select id, $1 from first_insert"

const GetURLsByUserID = "select id, user_id, original_url, short_url from urls t1, user_urls t2 where t1.id=t2.url_id and t2.user_id=$1"

const GetOriginalURLByShort = "select original_url from urls t1, user_urls t2 where t1.id =t2.url_id and t2.is_deleted=0 and t1.short_url=$1"

const GetOriginalURLByShortForUser = "select original_url from urls t1, user_urls t2 where t1.id=t2.url_id and t2.is_deleted=0 and t2.user_id=$1 and t1.short_url =$2"

const DeleteUserURL = "update user_urls t1 set is_deleted=1 from urls t2 where t1.url_id=t2.id and t1.user_id=$1 and t2.short_url=$2"
