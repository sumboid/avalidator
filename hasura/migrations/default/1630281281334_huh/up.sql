SET check_function_bodies = false;
CREATE FUNCTION public.set_current_timestamp_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$;
CREATE TABLE public.assignments (
    "user" uuid NOT NULL,
    task uuid NOT NULL,
    due_to time with time zone NOT NULL,
    id uuid DEFAULT public.gen_random_uuid() NOT NULL
);
CREATE TABLE public.cases (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    task uuid NOT NULL,
    input text NOT NULL,
    output text NOT NULL,
    created_by uuid NOT NULL
);
CREATE TABLE public.results (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    solution uuid NOT NULL,
    verdict text NOT NULL
);
CREATE TABLE public.solutions (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    assignment uuid NOT NULL,
    assets uuid NOT NULL
);
CREATE TABLE public.tasks (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    title text NOT NULL,
    description text NOT NULL
);
CREATE TABLE public.users (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    "group" text,
    role text NOT NULL
);
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignments_id_key UNIQUE (id);
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignments_pkey PRIMARY KEY ("user", task);
ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.solutions
    ADD CONSTRAINT solutions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
CREATE TRIGGER set_public_cases_updated_at BEFORE UPDATE ON public.cases FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_cases_updated_at ON public.cases IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_results_updated_at BEFORE UPDATE ON public.results FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_results_updated_at ON public.results IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_solutions_updated_at BEFORE UPDATE ON public.solutions FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_solutions_updated_at ON public.solutions IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_tasks_updated_at BEFORE UPDATE ON public.tasks FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_tasks_updated_at ON public.tasks IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_users_updated_at ON public.users IS 'trigger to set value of column "updated_at" to current timestamp on row update';
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignments_task_fkey FOREIGN KEY (task) REFERENCES public.tasks(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignments_user_fkey FOREIGN KEY ("user") REFERENCES public.users(id) ON UPDATE RESTRICT ON DELETE CASCADE;
ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON UPDATE RESTRICT;
ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_task_fkey FOREIGN KEY (task) REFERENCES public.tasks(id) ON UPDATE RESTRICT ON DELETE CASCADE;
ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_solution_fkey FOREIGN KEY (solution) REFERENCES public.solutions(id) ON UPDATE RESTRICT ON DELETE CASCADE;
ALTER TABLE ONLY public.solutions
    ADD CONSTRAINT solutions_assignment_fkey FOREIGN KEY (assignment) REFERENCES public.assignments(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
