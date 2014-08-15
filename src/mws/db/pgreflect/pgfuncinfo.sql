SELECT  proname, proargnames, pronargs,prosrc,format_type(prorettype, NULL),pg_get_function_identity_arguments(proname::regproc)
FROM    pg_catalog.pg_namespace n
JOIN    pg_catalog.pg_proc p
ON      pronamespace = n.oid
WHERE   nspname = 'public';



`SELECT
          proname,
          pronargs,
          proargmodes,
          array_agg(proargtypes), -- see here
           proallargtypes,
          proargnames,
          prodefaults,
          prorettype

FROM (
        SELECT
          p.proname,
          p.pronargs,
          p.proargmodes,
          format_type(unnest(p.proallargtypes), NULL) AS proargtypes, -- and here
          p.proallargtypes,
          p.proargnames,
          pg_get_expr(p.proargdefaults, 0) AS prodefaults,
          format_type(p.prorettype, NULL) AS prorettype
        FROM pg_catalog.pg_proc p
        JOIN pg_catalog.pg_namespace n
        ON n.oid = p.pronamespace
        WHERE n.nspname = 'public'
) x
GROUP BY proname, pronargs, proargmodes, proallargtypes, proargnames, prodefaults, prorettype`)
