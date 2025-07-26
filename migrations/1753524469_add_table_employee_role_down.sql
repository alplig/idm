\connect idm;

DROP TRIGGER IF EXISTS employee_set_timestamp ON public.employee;
DROP TRIGGER IF EXISTS role_set_timestamp ON public.role;
DROP FUNCTION IF EXISTS set_timestamp;

DROP TABLE IF EXISTS public.employee;
DROP TABLE IF EXISTS public.role;