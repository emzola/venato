CREATE TABLE IF NOT EXISTS project(
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    start_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    target_end_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    actual_end_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_on timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by bigint NOT NULL,
    modified_on timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    modified_by bigint NOT NULL,
    version integer NOT NULL DEFAULT 1
);