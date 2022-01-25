CREATE FUNCTION public.alter_entity(
    _entity_id INT,
    _code character varying
) RETURNS record
    LANGUAGE plpgsql
    AS $$
DECLARE
    -- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
    -- >> Description: Create/Edit entity                                           >>
    -- >> Version:     interview                                                    >>
    -- >> Date:        24/jan/2022                                                  >>
    -- >> Developer:   Edwin Plauchu Camacho                                        >>
    -- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

    current_moment timestamp with time zone := now();

    -- dump of errors
    rmsg text := '';

BEGIN

    CASE
        WHEN _entity_id = 0 THEN

            INSERT INTO entitys(
                code,
                last_touch_time,
                creation_time,
                blocked
            )VALUES(
                _code,
                current_moment,
                current_moment,
                false
            )RETURNING id INTO _entity_id;

        WHEN _entity_id > 0 THEN

            UPDATE entitys
            SET code = _code,
                last_touch_time = current_moment
            WHERE id = _entity_id;

        ELSE

            RAISE EXCEPTION 'negative entity identifier % is unsupported', _entity_id;


    END CASE;

    return ( _entity_id::integer, ''::text );

    EXCEPTION
        WHEN OTHERS THEN
            GET STACKED DIAGNOSTICS rmsg = MESSAGE_TEXT;
            return ( -1::integer, rmsg::text );

END;
$$;
