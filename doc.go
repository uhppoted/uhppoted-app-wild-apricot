// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppoted-app-wild-apricot integrates the uhppote-core API with a membership managed by Wild Apricot.

uhppoted-app-wild-apricot can be used from the command line but is really intended to be run from a cron job
to maintain the cards and permissions on a set of access controllers by deriving a unified access control list (ACL)
from a Wild Apricot membership list.

uhppoted-app-wild-apricot supports the following commands:

  - load-acl, to download an ACL from Wild Apricot to a set of access controllers
  - compare-acl, to compare an ACL from Wild Apricot with the cards and permissons on a set of access controllers
  - get-acl, to retrieve the ACL from Wild Apricot and store it to a TSV file
  - get-members, to retrieve a list of members from Wild Apricot and store it to a TSV file
  - get-groups, to retrieve a list of member groups from Wild Apricot and store it to a TSV file
  - get-doors, to extract a list of managed doors from the controller configuration file and store it to a TSV file
*/
package wild_apricot
