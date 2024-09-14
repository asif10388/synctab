create or replace function main.add_urls_v1(
	in _group_id uuid,
	in _user_id uuid,
	in _url text,
	in _title text,
	out _id uuid
) as $$
declare
begin
	insert into main.urls(group_id, user_id, url, title)
	values (_group_id, _user_id, _url, _title)
	
	returning id into _id;
end;
$$ language plpgsql;

/*SQLEND*/

create or replace function main.delete_url_v1(
	in _id uuid
) returns boolean as 
$$
declare
	_deleted boolean;
begin
	delete from main.urls
	where id = _id
	returning true into _deleted;
	
	return _deleted;
end;
$$ language plpgsql;

/*SQLEND*/