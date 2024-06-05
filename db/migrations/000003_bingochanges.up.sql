alter table public.bingos add column codephrase character varying(255) NOT NULL DEFAULT 'Bingo!';
alter table public.bingos add column ready boolean NOT NULL DEFAULT false;
