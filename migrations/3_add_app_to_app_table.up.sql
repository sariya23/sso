insert into app (app_id, name, secret)
values
(1, 'test', 'test-secret')
on conflict do nothing;