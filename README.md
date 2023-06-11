# ConfigBin

Service for storing configuration data of applications.

## Features

- The user can edit the configuration data through the web interface, simply by specifying a specific bin password.
- Third-party applications can receive configuration data with the simplest HTTP API, using the same password as the user for authorization. The password is sent via the HTTP header `X-ConfigBin-Password`.
- The configuration is encrypted with its password. The password is not stored in the database, so it cannot be recovered. The user has to save the password himself, otherwise he will not be able to restore the configuration.
- The bin ID is generated automatically as a UUID.
- TODO: There are code examples for different languages on the configuration data editing page.
- TODO: Configurations have versions and history. The user can revert to any version.
- TODO: The interface highlights syntax errors in the configuration, depending on the configuration format.
