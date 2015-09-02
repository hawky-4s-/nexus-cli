camunda-ci-helpers
==================

Usage
-----

Create repository:
POST https://app.camunda.com/nexus/service/local/repositories -> 201

{
  "data": {
    "repoType": "hosted",
    "id": "test-release",
    "name": "test-release",
    "writePolicy": "ALLOW_WRITE_ONCE",
    "browseable": true,
    "indexable": true,
    "exposed": true,
    "notFoundCacheTTL": 1440,
    "repoPolicy": "RELEASE",
    "provider": "maven2",
    "providerRole": "org.sonatype.nexus.proxy.repository.Repository",
    "downloadRemoteIndexes": false,
    "checksumPolicy": "IGNORE"
  }
}

{"data":{"contentResourceURI":"https://app.camunda.com/nexus/content/repositories/test-release","id":"test-release","name":"test-release","provider":"maven2","providerRole":"org.sonatype.nexus.proxy.repository.Repository","format":"maven2","repoType":"hosted","exposed":true,"writePolicy":"ALLOW_WRITE_ONCE","browseable":true,"indexable":true,"notFoundCacheTTL":1440,"repoPolicy":"RELEASE","downloadRemoteIndexes":false,"defaultLocalStorageUrl":"file:/home/java/jetty-nexus/sonatype-work/nexus/storage/test-release"}}

Delete repository:
DELETE https://app.camunda.com/nexus/service/local/repositories/test-release -> 204
