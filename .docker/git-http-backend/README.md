# Git-http-backend

This container uses Apache httpd to proxy request through git-http-backend. 
It specifically is used in the network to answer git request to github.com for
any of its remote needs.

All test that make remote request to github.com will
land here instead of the real github.com. Pretty neat. Not only that, it will
correctly answer those request in a way that the git CLI (as a client) will
understand. And since git http-backend is a legit and a great way to
mock response to the git CLI client. Even the repo are real repositories, but
those are manually made to suit the test. Sorry for all the extra reading, maybe
only the first paragraph was necessary.

1. First learn how to configure Apache. Not all of it, but enough to get the job
done and well enough to explain what you did and why. Use information on the
[httpd] container README as a guide.
2. Configure Apache [httpd.conf] with the following settings:
    1. Search, find, and remove the `#` at the beginning of the line to enable
       the following modules:
       * mod_cgi
       * mod_alias
       * mod_env
       * mpm_prefork_module
       * cgi_module
       * mod_authn_core
       * mod_authz_core
   NOTE: To disable something simply add the `#` at the beginning of a line.
    2. Disable `mpm_event_module` because 2 mpm's cannot be enabled at the same
       time. `mpm_prefork_module` seems required for CGI.
    3. Configure a virtual host to map to the directory where we cloned all the
       repositories to the site. You can add `github.com` as a VirtualHost to
       the [httpd.conf] custom Apache config, like so:
       ```httpdconf
       <VirtualHost *:443>
           DocumentRoot "/usr/local/apache2/htdocs/kohirens"
           ServerName github.com
    
           SetEnv GIT_PROJECT_ROOT /usr/local/apache2/htdocs/kohirens
           SetEnv GIT_HTTP_EXPORT_ALL 1
           ScriptAlias /kohirens/ /usr/libexec/git-core/git-http-backend/
    
           # This is not strictly necessary using Apache and a modern version of
           # git-http-backend, as the webserver will pass along the header in the
           # environment as HTTP_GIT_PROTOCOL, and http-backend will copy that into
           # GIT_PROTOCOL. But you may need this line (or something similar if you
           # are using a different webserver), or if you want to support older Git
           # versions that did not do that copying.
           #
           # Having the webserver set up GIT_PROTOCOL is perfectly fine even with
           # modern versions (and will take precedence over HTTP_GIT_PROTOCOL,
           # which means it can be used to override the client's request).
           SetEnvIf Git-Protocol ".*" GIT_PROTOCOL=$0
    
           RewriteCond %{QUERY_STRING} service=git-receive-pack [OR]
           RewriteCond %{REQUEST_URI} /git-receive-pack$
           RewriteRule ^/kohirens/ - [E=AUTHREQUIRED:yes]
    
           <LocationMatch "^/kohirens/">
               Order Deny,Allow
               Deny from env=AUTHREQUIRED
    
               AuthType Basic
               AuthName "Git Access"
               AuthBasicProvider file
               AuthUserFile "/usr/local/apache2/conf/passwords"
               Require user x-access-token
           </LocationMatch>
       </VirtualHost>
       ```
3. The `avr` application will usa a GitHub fine-grained access token when
   it makes request to push changes. That token along with a static username
   will be passed to Apache, which then will use HTTP basic authorization to
   authenticate any such request. See the Apache httpd
   [Authentication and Authorization] docs to learn how to set that up. Create
   the password file as explained in the docs, like so:
   ```shell
   htpasswd -c /usr/local/apache/passwd/passwords x-access-token
   ```
   Once you make the file, you'll need to copy it out of the container, then add
   it to the image when building.
   NOTE: The VirtualHost make reference to this file, so be careful to place it
         where the VirtualHost will look for it, otherwise test will fail.
4. git-http-backend will need real git repositories to send back responses to
   the Git CLI, so those are stored as `*.bundle` files in this repository. When
   the image is built they are copied into it and then a shell script
   clones them into the `/usr/local/apache2/htdocs/kohirens` directory inside
   the image. This directory serves aa the Apache webserver `DocumentRoot` and
   also the `GIT_PROJECT_ROOT`. Test are made to request repositories that exist
   in this directory.
5. git-http-backend also requires SSL/TlS/HTTPS to communicate successfully. You
   may get a 403 if you try to use HTTP. You can generate an SSL cert/key with
   the [gen-ss-cert.sh] script provided in this repository. You
   need to do this before you build the image, as the build requires it to be
   at a specific location. It is also added to the web containers Root-CA to
   allow successful TLS verification for request from that container. Place
   the generated files as `.docker/ssl/certs/server.pem` and
   `.docker/ssl/private/server.key`. Building the image will take care of the
   rest.
   NOTE: The SSL cert/key pair are copied into the container where Apache
   expects to find them see "SSLCertificateFile" and "SSLCertificateKeyFile"
   settings in the [httpd-ssl.conf] file.

## References

---

[httpd]: https://hub.docker.com/_/httpd
[Authentication and Authorization]: https://httpd.apache.org/docs/2.4/howto/auth.html
[httpd.conf]: /.docker/git-http-backend/my-httpd.conf
[httpd-ssl.conf]: /.docker/git-http-backend/my-httpd-ssl.conf
[gen-ss-cert.sh]: /.docker/git-http-backend/gen-ss-cert.sh
