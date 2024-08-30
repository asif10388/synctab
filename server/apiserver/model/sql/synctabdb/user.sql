create or replace function main.add_user_v1(
	in _email text,
	in _username text,
	in _password text,
	out _id uuid
) as $$
declare
begin
	if exists(select 1 from main.users where email = _email) then
		raise exception 'user already exists';
	end if;

	insert into main.users(email, username, passwordhash)
	values (_email, _username, _password)

	returning id into _id;
end;
$$ language plpgsql;

/*SQLEND*/