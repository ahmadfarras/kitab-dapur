CREATE TABLE kitabs (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE recipes (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    kitab_id    uuid NOT NULL REFERENCES kitabs (id) ON DELETE CASCADE,
    title       text NOT NULL,
    author      text NOT NULL,
    ingredients text[] NOT NULL,
    steps       text[] NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

-- Serves "a kitab's recipes, newest first".
CREATE INDEX recipes_kitab_id_created_at ON recipes (kitab_id, created_at DESC, id DESC);
