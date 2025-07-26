\connect idm;

CREATE TABLE public.role
(
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.employee
(
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION set_timestamp()
    RETURNS trigger AS
$$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER employee_set_timestamp
    BEFORE UPDATE
    ON public.employee
    FOR EACH ROW
    EXECUTE FUNCTION set_timestamp();

CREATE TRIGGER role_set_timestamp
    BEFORE UPDATE
    ON public.role
    FOR EACH ROW
    EXECUTE FUNCTION set_timestamp();