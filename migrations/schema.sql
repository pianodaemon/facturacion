--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3
-- Dumped by pg_dump version 13.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

alter database interview SET TIME ZONE 'America/Mexico_City';

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apps; Type: TABLE; Schema: public; Owner: postgres
--


CREATE TABLE public.entitys (
    id serial NOT NULL,
    code character varying NOT NULL,
    last_touch_time timestamp with time zone NOT NULL,
    creation_time timestamp with time zone NOT NULL,
    blocked boolean DEFAULT false NOT NULL
);


--
-- Name: entitys entitys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entitys
    ADD CONSTRAINT entitys_pkey PRIMARY KEY (id);
