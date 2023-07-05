# Certificate

## Run
* To bring up all components, run `docker-compose up`
* To run unit tests, run `docker-compose -f docker-compose-test.yml up`

## API Endpoints
* `POST /user`
  * Takes in JSON fields `name`, `email`, `password`
  * Add a new user with above attributes if there's no user with the same email
  * Returns the user with its newly generated UUID
* `GET /user`
  * Takes in a JSON field `uuid`
  * Returns all the user's attributes if found and not deleted
* `DELETE /user`
  * Takes in a JSON field `uuid`
  * Marks the user as inactive and appear to be deleted in subsequent requests
* `POST /cert`
  * Takes in JSON fields `user_uuid`, `private_key`, `body`
  * Add them as a new certificate
  * Returns the certificate with its newly generated UUID
* `GET /cert`
  * Takes in a JSON field `user_uuid`
  * Returns a list of active certificates belonging to `user_uuid`
* `PATCH /cert`
  * Take in JSON fields `uuid`, `user_uuid` and `active`
  * Deactivate/activate the certificate according to `active`
  * Posts notification to `ENDPOINT` env, set to https://enczcbi39ybms.x.pipedream.net/ in this repo
  * Returns error and does not notify if cert is already active / already inactive

## Assumptions
### Certificate activation/deactivation notifications
* Creating a new certificate counts as activating it, so creation warrants a POST to our http bin
* Activation/deactivation messages to HTTP bin can tolerate some delay
* Occasional duplicate activation/deactivate messages to HTTP bin are not a big problem
### User deletion
* User deletion is implemented as deactivation, we do not want to immediately lose all user data upon deletion
  * API behaviors after deactivation simulates deletion - i.e. trying to get a deactivate user returns error
* User's certificates do not have to be deactivated upon user deletion

## Out of Scope because Out Of Time
* Specific non-200 HTTP status codes, using 500 for everything
* Config and credentials: using ENVs and hard coded values
* Input validation for APIs and libraries
* Pagination on the list of certificates
* Integration testing
* Unit testing for `notifier` service, and `certificate/router` package