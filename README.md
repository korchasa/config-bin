# ConfigBin

A service for remote storage and editing of configurations. The user edits the configuration via a web interface. The application loads the configuration through a simple API.

## Features

- The user can edit the configuration data through the web interface, simply by specifying the specific bin password.
- Third-party applications can receive configuration data through the simplest HTTP API, using the same password as the user for authorization. The password is sent via the HTTP Basic auth(with any username or without it).
- The configuration is encrypted using its password. The password is not stored in the database, so it cannot be recovered. The user needs to save the password themselves, otherwise, they will not be able to restore the configuration.
- The bin ID is automatically generated UUID.
- TODO: Code examples in different languages are available on the configuration data editing page.
- TODO: Configurations have version histories. The user can revert to any version.
- TODO: The interface highlights syntax errors in the configuration, depending on the configuration format.
