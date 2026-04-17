# Changelog

## 2.0.0-alpha.5 (2026-04-17)

Full Changelog: [v2.0.0-alpha.4...v2.0.0-alpha.5](https://github.com/G-Core/terraform-provider-gcore/compare/v2.0.0-alpha.4...v2.0.0-alpha.5)

### Chores

* **tests:** bump steady to v0.22.1 ([7622522](https://github.com/G-Core/terraform-provider-gcore/commit/76225229a9f0cbe7afa1f2464593b3893cbf31c7))

## 2.0.0-alpha.4 (2026-04-17)

Full Changelog: [v2.0.0-alpha.3...v2.0.0-alpha.4](https://github.com/G-Core/terraform-provider-gcore/compare/v2.0.0-alpha.3...v2.0.0-alpha.4)

### Features

* add per-resource api permissions to schema description ([fe2ef35](https://github.com/G-Core/terraform-provider-gcore/commit/fe2ef3539059262ec0f5ef34d3430b97e2cd9318))
* **api:** aggregated API specs update ([3a7edbb](https://github.com/G-Core/terraform-provider-gcore/commit/3a7edbb756a95c4cde4121c142343ef53dbce400))
* **api:** aggregated API specs update ([5cca998](https://github.com/G-Core/terraform-provider-gcore/commit/5cca9989fb646cca3fe2add7135fd9ada16b25ee))
* **api:** aggregated API specs update ([e29876a](https://github.com/G-Core/terraform-provider-gcore/commit/e29876a8e8ce79aa3c8260de0c935609d358254d))
* **api:** aggregated API specs update ([5e9437f](https://github.com/G-Core/terraform-provider-gcore/commit/5e9437f8ff350b375fde3adedb8b827794641ff0))
* **api:** aggregated API specs update ([88f1919](https://github.com/G-Core/terraform-provider-gcore/commit/88f1919fd51eaa4930d99e6791ffd349ac2a9484))
* **api:** aggregated API specs update ([a3837f4](https://github.com/G-Core/terraform-provider-gcore/commit/a3837f4d7d7a075b095e82f1753787cddd32e298))
* **api:** aggregated API specs update ([8c9baaa](https://github.com/G-Core/terraform-provider-gcore/commit/8c9baaa25b48963842f5e21bdcb9d9bbf9aa8fbb))
* **api:** aggregated API specs update ([16692ef](https://github.com/G-Core/terraform-provider-gcore/commit/16692ef28abb3f86dd934084e3b014c7377a091c))
* **api:** aggregated API specs update ([2fd8a29](https://github.com/G-Core/terraform-provider-gcore/commit/2fd8a2920732b014c98c7e73cf3f2bfa18d130b6))
* **api:** aggregated API specs update ([3187470](https://github.com/G-Core/terraform-provider-gcore/commit/3187470a92e1efc60b376beda978f278587cc13e))
* **api:** aggregated API specs update ([1837170](https://github.com/G-Core/terraform-provider-gcore/commit/18371707f6f1bb4544f9bdcd6c1c529e29c12235))
* **api:** api update ([4d81d0a](https://github.com/G-Core/terraform-provider-gcore/commit/4d81d0aa7c53aab279875100c87fb360b07128d1))
* **api:** feat(cli): add x-stainless-skip for limit/offset parameters in CLI target ([d34a7e0](https://github.com/G-Core/terraform-provider-gcore/commit/d34a7e0983268ed3b48a694e392e8e5aaa31d16c))
* **api:** fix: skip offset parameter only for cli, not terraform ([f35645a](https://github.com/G-Core/terraform-provider-gcore/commit/f35645a1cbd67eefc73dc32f24d2b05059025908))
* **api:** fix(cdn): harmonize pagination across CDN list endpoints ([29cd2d9](https://github.com/G-Core/terraform-provider-gcore/commit/29cd2d9a81b4e29f98d15c0538757d2316baf52f))
* **cdn:** add S3 inline origins for origin groups ([aee1680](https://github.com/G-Core/terraform-provider-gcore/commit/aee1680c09f9a6eae1cb5f00026c5c656a929bf7))
* **cloud:** add find_one_by to cloud_region data source ([7d3fcbd](https://github.com/G-Core/terraform-provider-gcore/commit/7d3fcbdbdd40f70c04544cf8d604c9bd8cced13b))


### Bug Fixes

* align cdn_resource_rule data source model tag with schema ([2ac73c9](https://github.com/G-Core/terraform-provider-gcore/commit/2ac73c9777bf796899e2cc512d2ad8e6b8643ffe))
* **cdn:** deactivate cdn_resource_rule before deletion, mark deleted as no_refresh ([ec83f2d](https://github.com/G-Core/terraform-provider-gcore/commit/ec83f2dcec4d2b59a5cbddca0e5e86cb4c5469cb))
* **ci:** add disk cleanup before goreleaser ([976ab43](https://github.com/G-Core/terraform-provider-gcore/commit/976ab4359ad552508c52c300b61ba459489d7ac9))
* **ci:** in custom setup-go, pass through go-version and cache-dependency-path ([a6970f9](https://github.com/G-Core/terraform-provider-gcore/commit/a6970f90830450f5ce798ef985a159050d75035d))
* **cloud:** allow tags to be refreshed ([44b8965](https://github.com/G-Core/terraform-provider-gcore/commit/44b8965ce20afb896aa331036b89ec0c7bcf5166))
* **cloud:** prevent forced replacement of project_id/region_id after import ([1e85cb9](https://github.com/G-Core/terraform-provider-gcore/commit/1e85cb90e6ef60533463b3562d4821f58e439f64))
* **cloud:** restore cloud_k8s_cluster_kubeconfig data source ([032b755](https://github.com/G-Core/terraform-provider-gcore/commit/032b75553ed30bd3c3de23e9f160321a30d460d8))
* **cloud:** restore missing customfield import in cloud_instance_image ([3e1c5b2](https://github.com/G-Core/terraform-provider-gcore/commit/3e1c5b2ad6018155e76dcfd47811df708c29c7af))
* **cloud:** restore params for inference flavors list data source ([70b0109](https://github.com/G-Core/terraform-provider-gcore/commit/70b010998820fa1e161a07d1bacd172daf060b27))
* fall back to main branch if linking fails in CI ([1f6c142](https://github.com/G-Core/terraform-provider-gcore/commit/1f6c1421e102af0613623c67aaed045a13f5e449))
* fix for failing to drop invalid module replace in link script ([add67cf](https://github.com/G-Core/terraform-provider-gcore/commit/add67cf831de25a5d926c844f70235db211c9dbe))
* fix quoting typo ([ccf4b27](https://github.com/G-Core/terraform-provider-gcore/commit/ccf4b27325110b4bbc1cbcb5b9a31500419c4b3f))
* remove jsonencode() requirement for content fields ([f6b926d](https://github.com/G-Core/terraform-provider-gcore/commit/f6b926df5438bdd59f1286cd4cf333cf64ff9c53))
* update custom code for Go SDK v0.41.0 compatibility ([ec863f1](https://github.com/G-Core/terraform-provider-gcore/commit/ec863f156372b899a704a1e015284a3123b74880))


### Chores

* add comment ([1b19948](https://github.com/G-Core/terraform-provider-gcore/commit/1b1994840da9922d53265c2ae97911e5b64ebe05))
* **internal:** codegen related update ([75911d6](https://github.com/G-Core/terraform-provider-gcore/commit/75911d66de2abcb7b69746a819499753a5897969))
* pin go releaser version ([5fc2d04](https://github.com/G-Core/terraform-provider-gcore/commit/5fc2d046cbf1200fe9bc1e7abb34d0219d18b36c))
* regenerate docs ([fd0fe40](https://github.com/G-Core/terraform-provider-gcore/commit/fd0fe409afdc40ca8e78883d10bd048522fe91fc))
* skip shared CDN limit $ref parameter in terraform codegen ([2ccac86](https://github.com/G-Core/terraform-provider-gcore/commit/2ccac86f61c7e2aeac8fb003b2554eb5fde04f18))
* skip shared CDN limit $ref parameter in terraform codegen ([3f9e63c](https://github.com/G-Core/terraform-provider-gcore/commit/3f9e63cf26a2cb8f4afcccc1e2a01559cb16c786))
* **tests:** bump steady to v0.20.1 ([5126625](https://github.com/G-Core/terraform-provider-gcore/commit/5126625b77f16ace6b13d5482f7bf62d581dc116))
* **tests:** bump steady to v0.20.2 ([d17bfc7](https://github.com/G-Core/terraform-provider-gcore/commit/d17bfc7390fc0c67218c82c0b900ad3b81d41fe0))


### Documentation

* clarify provider-level project/region IDs as dev convenience ([9f148a5](https://github.com/G-Core/terraform-provider-gcore/commit/9f148a5d3c4df1ff2bd0dc573c0cbe12a0c13c50))
* update examples ([8562e89](https://github.com/G-Core/terraform-provider-gcore/commit/8562e891578446bffc2a2ee53afb7add7148ab79))

## 2.0.0-alpha.3 (2026-03-30)

Full Changelog: [v2.0.0-alpha.2...v2.0.0-alpha.3](https://github.com/G-Core/terraform-provider-gcore/compare/v2.0.0-alpha.2...v2.0.0-alpha.3)

### Features

* add release workflow skill for provider releases ([2a77358](https://github.com/G-Core/terraform-provider-gcore/commit/2a77358539eceb6a5e79c0559584ecea1e7884e6))
* **api:** aggregated API specs update ([2d16176](https://github.com/G-Core/terraform-provider-gcore/commit/2d161766414a196cb585e62e6f45fe569d5a1eb9))
* **api:** aggregated API specs update ([2c398a8](https://github.com/G-Core/terraform-provider-gcore/commit/2c398a800c88221b2df00fa17030d12fa8557577))
* **api:** aggregated API specs update ([0a61dc7](https://github.com/G-Core/terraform-provider-gcore/commit/0a61dc7cdad7e5a2af7587c1050eab9e80a56471))
* **api:** feat(cdn): enable terraform support for cdn_rule_template ([f492d9a](https://github.com/G-Core/terraform-provider-gcore/commit/f492d9a63dd9139365eb1044595aed4c5d227a05))
* **cdn:** enable terraform support for cdn_rule_template ([3fc2b91](https://github.com/G-Core/terraform-provider-gcore/commit/3fc2b91ef2d1f1c0e978370488c9f50961383bfe))
* **cloud:** add support for Baremetal servers in Terraform ([0e5d1a3](https://github.com/G-Core/terraform-provider-gcore/commit/0e5d1a388c102723511f990f952b668ada8f7ee7))
* **cloud:** add support for Baremetal servers in Terraform ([dc2d54f](https://github.com/G-Core/terraform-provider-gcore/commit/dc2d54fcdf936f0d47839c802b7db88d13650935))
* **cloud:** enable Terraform data source for inference flavors ([6f6d5cf](https://github.com/G-Core/terraform-provider-gcore/commit/6f6d5cf7437117941415b8cfc9e0e1af7b057d39))
* **fastedge:** enable fastedge_template for Terraform ([ae3b638](https://github.com/G-Core/terraform-provider-gcore/commit/ae3b638234493f97318a52e835b805b47a020181))
* **fastedge:** enable fastedge_template for Terraform ([f5fe838](https://github.com/G-Core/terraform-provider-gcore/commit/f5fe8383ae827a67f1adc30574d86ffeb09c339f))


### Bug Fixes

* **cloud:** harden GPU cluster credentials handling ([b4ba5dc](https://github.com/G-Core/terraform-provider-gcore/commit/b4ba5dc2140276fc1d0dee0639169229432817c1))
* **cloud:** remove client_id from cloud project data sources ([1504311](https://github.com/G-Core/terraform-provider-gcore/commit/1504311863d34a21a5959df8f87ea94df3acf734))
* **cloud:** remove ddos_profile from resource and data sources of gcore_cloud_load_balancer ([753bff2](https://github.com/G-Core/terraform-provider-gcore/commit/753bff2e5a2d2b140a1392a7bb8e7465a2eafe42))
* **cloud:** remove dead planNames variable in plan modifier of cloud_k8s_cluster ([c732082](https://github.com/G-Core/terraform-provider-gcore/commit/c7320825053eb7d3716449352a53bd47eaaf0ee9))
* **fastedge:** normalize ImportState response to prevent post-import drift ([453ba77](https://github.com/G-Core/terraform-provider-gcore/commit/453ba778f0db443196c7085d7308c07e71aacedb))
* patch style requests should never send empty json body for objects ([8053c46](https://github.com/G-Core/terraform-provider-gcore/commit/8053c465f5a9e81ce8dbde02bffbf4103f5ebed6))


### Chores

* **internal:** update multipart form array serialization ([fec341f](https://github.com/G-Core/terraform-provider-gcore/commit/fec341face73433749aace59fadcafd2018ae13b))
* **tests:** bump steady to v0.19.6 ([d8f4fc1](https://github.com/G-Core/terraform-provider-gcore/commit/d8f4fc1e5c47ab38e433bb04fb78c911821c88e5))
* **tests:** bump steady to v0.19.7 ([50456eb](https://github.com/G-Core/terraform-provider-gcore/commit/50456eb6f82c4241415943fe62bc301fe25536be))


### Documentation

* update provider version example to 2.0.0-alpha.3 ([3bdb94f](https://github.com/G-Core/terraform-provider-gcore/commit/3bdb94f144c2ec62a1d72af9160560733bd74014))

## 2.0.0-alpha.2 (2026-03-23)

Full Changelog: [v2.0.0-alpha.1...v2.0.0-alpha.2](https://github.com/G-Core/terraform-provider-gcore/compare/v2.0.0-alpha.1...v2.0.0-alpha.2)

### Features

* **api:** aggregated API specs update ([d39c234](https://github.com/G-Core/terraform-provider-gcore/commit/d39c234101c273a9b64de6fd98abd1efc758f663))
* **api:** aggregated API specs update ([bd6b4b5](https://github.com/G-Core/terraform-provider-gcore/commit/bd6b4b5284ce172e6ff00113b07d550fe3c860bd))
* **api:** aggregated API specs update ([f0c436d](https://github.com/G-Core/terraform-provider-gcore/commit/f0c436ddb8c2dd00fe158548e20fdfcb554ca004))
* **api:** aggregated API specs update ([3038dfb](https://github.com/G-Core/terraform-provider-gcore/commit/3038dfb4d47b4fc72a63bc68f927777ff218b934))
* **api:** feat(cdn): enable terraform support for cdn_rule_template ([2b689b0](https://github.com/G-Core/terraform-provider-gcore/commit/2b689b08978d6f5964551e079f757c089652f1f5))
* **api:** feat(cloud): add support for GPU baremetal clusters in Terraform ([b6edc7c](https://github.com/G-Core/terraform-provider-gcore/commit/b6edc7ccdb8b6fbcd5bc02bfc5e3f40425c7a671))
* **api:** Revert "feat(cdn): enable terraform support for cdn_rule_template" ([46bc3ca](https://github.com/G-Core/terraform-provider-gcore/commit/46bc3ca531e9d60a84f82832c28cb57aef04cc01))
* **cdn:** enable terraform support for cdn_rule_template ([303a9c4](https://github.com/G-Core/terraform-provider-gcore/commit/303a9c4ad2d79794ab680a0fbd509b5640077c30))
* **cloud:** add gcore_cloud_k8s_kubeconfig data source ([d003c35](https://github.com/G-Core/terraform-provider-gcore/commit/d003c35fccdee0646c434bbb51efe405599c7f91))
* **cloud:** add support for cloud_projects resource in Terraform ([4b35d5a](https://github.com/G-Core/terraform-provider-gcore/commit/4b35d5a5e230bde3b310d69145dac27de76a06c4))
* **cloud:** add support for cloud_projects resource in Terraform ([21a2be4](https://github.com/G-Core/terraform-provider-gcore/commit/21a2be46834b8ead0916ce718b773ed4d621d0b1))
* **cloud:** add support for GPU baremetal clusters in Terraform ([87b7adb](https://github.com/G-Core/terraform-provider-gcore/commit/87b7adb180370d4dd8290d477458f94268d338e1))
* **cloud:** enable Terraform data source for GPU baremetal clusters ([1612a42](https://github.com/G-Core/terraform-provider-gcore/commit/1612a42875494dfe15e3ffdc27b663f1cafbfa86))
* **cloud:** GPU baremetal cluster full CRUD and update support ([6a9b527](https://github.com/G-Core/terraform-provider-gcore/commit/6a9b527ff8a8c5b3ee463fef7e6e7132163ae666))
* **cloud:** remove create_router field from gcore_cloud_network ([87d1dde](https://github.com/G-Core/terraform-provider-gcore/commit/87d1dde603f3e3eea2aab28564c420f9ee6bcac5))
* **cloud:** remove deprecated timeout_client_data from gcore_cloud_load_balancer_pool ([5a8385d](https://github.com/G-Core/terraform-provider-gcore/commit/5a8385d796a27b6646305bda1cc68cf0d0abdb3e))
* **cloud:** remove deprecated timeout_member_connect and timeout_member_data from gcore_cloud_load_balancer_listener ([99f9be3](https://github.com/G-Core/terraform-provider-gcore/commit/99f9be3567006ae15e11a23a53ce580f8711928c))
* **cloud:** remove reservation block from gcore_cloud_reserved_fixed_ip ([832fbc3](https://github.com/G-Core/terraform-provider-gcore/commit/832fbc3c4cdf1ce7d409488b182b938c54ce6251))
* **cloud:** remove router attachment fields from gcore_cloud_network_subnet ([833ec8c](https://github.com/G-Core/terraform-provider-gcore/commit/833ec8c7776449e4dcc5b61969025126d8102bb8))


### Bug Fixes

* **ci:** gracefully skip SDK link when branch doesn't exist ([b00d880](https://github.com/G-Core/terraform-provider-gcore/commit/b00d880abe38458d5f4e8a9533308bd4ac800484))
* **cloud:** handle unknown values in GPU virtual cluster validators ([4069780](https://github.com/G-Core/terraform-provider-gcore/commit/4069780a5b73221b69e8957dc3468a9b7c0c6de7))
* **cloud:** infer Type field on import to prevent plan diff ([ddb1d7e](https://github.com/G-Core/terraform-provider-gcore/commit/ddb1d7e232ea060319e35996ef72e73cb198ee53))
* **cloud:** support PATCHing tags on gcore_cloud_instance update ([e989285](https://github.com/G-Core/terraform-provider-gcore/commit/e98928575813b5c7875b3b78513c489b8cf218bd))
* **cloud:** use DeleteAndPoll to wait for resource removal ([5539e98](https://github.com/G-Core/terraform-provider-gcore/commit/5539e98efb9abe8b9dfe5b0a2a21a12466f5e15f))
* **cloud:** use Update method for GPU virtual cluster tags ([cf80999](https://github.com/G-Core/terraform-provider-gcore/commit/cf8099979b70ff7cb79f4c3d19fe5411e2540682))
* **cloud:** use UseStateForUnknown for size of cloud_gpu_virtual_cluster_image ([49ba618](https://github.com/G-Core/terraform-provider-gcore/commit/49ba61872c7270c3f3d4656a8c5db5b2f21d927d))
* **dns:** fix update drift, name immutability, add tests for dns_network_mapping ([b08e721](https://github.com/G-Core/terraform-provider-gcore/commit/b08e721eec561a4698aafa4f21cc76ef2a3e9c05))
* **fastedge:** remove id=0 workarounds after API fix ([789a674](https://github.com/G-Core/terraform-provider-gcore/commit/789a6749e53908aef69273dd735f7d03f1056cef))
* **fastedge:** rename path parameters to match updated OpenAPI spec ([3a2115f](https://github.com/G-Core/terraform-provider-gcore/commit/3a2115f6de2409ce3437822e9e292b012efe6dc7))
* **fastedge:** rename path parameters to match updated OpenAPI spec ([cb03da1](https://github.com/G-Core/terraform-provider-gcore/commit/cb03da139183caa25e3a1b9665dae2c57854a237))
* **fastedge:** skip ValidateConfig when binary/template is unknown ([4d459c6](https://github.com/G-Core/terraform-provider-gcore/commit/4d459c61c16ccf14afb5a109ab22c5c63959fd47))
* improve linking behavior when developing on a branch not in the Go SDK ([8d832d5](https://github.com/G-Core/terraform-provider-gcore/commit/8d832d5774840e2b84ef1c18179c2120b347f951))
* improved workflow for developing on branches ([e500cee](https://github.com/G-Core/terraform-provider-gcore/commit/e500ceea19ab7a8a5e78520e4a5534eea316ef55))
* **internal:** fix issue with some nested custom attributes ([8338de4](https://github.com/G-Core/terraform-provider-gcore/commit/8338de454f771800ea61ef742bcf7c32f5c368fa))
* no longer require an API key when building on production repos ([6c3fce5](https://github.com/G-Core/terraform-provider-gcore/commit/6c3fce51b265cfb59558220b286cc4d96a98ba19))


### Reverts

* enable fossa failing on error ([351620c](https://github.com/G-Core/terraform-provider-gcore/commit/351620ca172eb16f31e5ef1fdb91eee19dfe2e3a))


### Chores

* add missing cdn_resource override name in config ([1011ce7](https://github.com/G-Core/terraform-provider-gcore/commit/1011ce7f2956dcabcaffe8af5b47e7d1f371c006))
* add missing cdn_resource override name in config ([b32406a](https://github.com/G-Core/terraform-provider-gcore/commit/b32406af02bb78935ba35ae145258a3a8adcbed8))
* **docs:** update terraform-plugin-docs to v0.24.0 ([21d566a](https://github.com/G-Core/terraform-provider-gcore/commit/21d566a256273a1dd6045d9e61c15fe32e35073a))
* **internal:** codegen related update ([b7bd73d](https://github.com/G-Core/terraform-provider-gcore/commit/b7bd73d63540b89d2aeb1a75c64befde4b73e7f7))
* **internal:** tweak CI branches ([1b6aedf](https://github.com/G-Core/terraform-provider-gcore/commit/1b6aedf09031fe2caee7ace6e3feee49ab287875))
* **internal:** update gitignore ([c01a899](https://github.com/G-Core/terraform-provider-gcore/commit/c01a8994da3a5bc94b525f4da598a8008d73bdc0))
* **tests:** bump steady to v0.19.4 ([722d95f](https://github.com/G-Core/terraform-provider-gcore/commit/722d95f7cde269727344debdf2a8a6f31a359c3b))
* **tests:** bump steady to v0.19.5 ([9da6050](https://github.com/G-Core/terraform-provider-gcore/commit/9da605055603d5b76cf5cf1f68f3b309234a680b))


### Documentation

* add missing import examples for 3 resources ([1603c04](https://github.com/G-Core/terraform-provider-gcore/commit/1603c04dbe02978a2f0e5974ac5f0d27a97b6bf5))
* add support contact disclaimer to provider index ([881b03f](https://github.com/G-Core/terraform-provider-gcore/commit/881b03f45cb78b796fb4c25b67263a52646f9593))
* regenerate index.md from template and fix instance examples ([a6be214](https://github.com/G-Core/terraform-provider-gcore/commit/a6be214c7061b0145bba482feeb5a3f91be9b160))


### Refactors

* **tests:** switch from prism to steady ([699f2fb](https://github.com/G-Core/terraform-provider-gcore/commit/699f2fbb086b94e1834c5e434e96a7b5e13a5b86))

## 2.0.0-alpha.1 (2026-03-09)

Full Changelog: [v0.0.1-alpha.0...v2.0.0-alpha.1](https://github.com/G-Core/terraform-provider-gcore/compare/v0.0.1-alpha.0...v2.0.0-alpha.1)

### ⚠ BREAKING CHANGES

* **api:** change type casing from Cdn* to CDN*
* **cloud:** use create and update v2 endpoints for security groups
* **cloud:** use v2 endpoint for floating IPs updates
* change naming for POST, PUT, PATCH, DELETE models
* **cloud:** restructure to be gpu_virtual.clusters
* **cloud:** replace PUT /cloud/v1/l7policies with PATCH
* **cloud:** rename load balancer pool member methods to create/delete
* **cloud:** k8s references from k8 to k8s

### Features

* add acctest and configurability subagents for opencode ([a519fe5](https://github.com/G-Core/terraform-provider-gcore/commit/a519fe58c09e0089335bf136f41ce58b5c8585b8))
* add acctest skill ([6ad59db](https://github.com/G-Core/terraform-provider-gcore/commit/6ad59db958590dabb3201461cce4a7f8a0de258a))
* add conditional clearing for external_gateway_info to reduce drift ([da4792a](https://github.com/G-Core/terraform-provider-gcore/commit/da4792a911903d4a87f542c554fe9a43c9486746))
* add descriptions for all Terraform-enabled resources ([eec7e06](https://github.com/G-Core/terraform-provider-gcore/commit/eec7e06fd73db2af77fe2431cbc50fa69450c5f3))
* add lefthook pre-commit hooks ([#100](https://github.com/G-Core/terraform-provider-gcore/issues/100)) ([b42a9a7](https://github.com/G-Core/terraform-provider-gcore/commit/b42a9a78e6adc13ff995d7afb6305369a00c84c9))
* add plan-modifiers skill ([fd22c34](https://github.com/G-Core/terraform-provider-gcore/commit/fd22c346cd193fa58996f8ccfd63784115bb72b4))
* **api:** aggregated API specs update ([a1a2d6a](https://github.com/G-Core/terraform-provider-gcore/commit/a1a2d6a6cd704b089450c36c3cbd3004278eb67d))
* **api:** aggregated API specs update ([7f47b1d](https://github.com/G-Core/terraform-provider-gcore/commit/7f47b1db759db41386877868a7ffa78a9545e286))
* **api:** aggregated API specs update ([38e6812](https://github.com/G-Core/terraform-provider-gcore/commit/38e68129413ce791443cf9bbd97f4fc127036718))
* **api:** aggregated API specs update ([1a469f8](https://github.com/G-Core/terraform-provider-gcore/commit/1a469f894bc2c3fc7b8aefd762a581e9e5ae2432))
* **api:** aggregated API specs update ([bd16fae](https://github.com/G-Core/terraform-provider-gcore/commit/bd16fae1b618c2d45a133045b63863772879c07f))
* **api:** aggregated API specs update ([a95b042](https://github.com/G-Core/terraform-provider-gcore/commit/a95b04288838451513b65402ab19d96bee2fe727))
* **api:** aggregated API specs update ([11d6170](https://github.com/G-Core/terraform-provider-gcore/commit/11d61703d98f4ef867cb2bbc21ab4aa7c48cc559))
* **api:** aggregated API specs update ([0379380](https://github.com/G-Core/terraform-provider-gcore/commit/037938037a2c6dd2dba72a774e5909cc80de57ce))
* **api:** aggregated API specs update ([c6091d0](https://github.com/G-Core/terraform-provider-gcore/commit/c6091d0aa9434f18462a7d8e228a8ad4873f0853))
* **api:** aggregated API specs update ([8460e9f](https://github.com/G-Core/terraform-provider-gcore/commit/8460e9fffaeed49c759285639d1436f799a3568a))
* **api:** aggregated API specs update ([123bf47](https://github.com/G-Core/terraform-provider-gcore/commit/123bf4769d4b160288db23cef3e7d5f9858f5704))
* **api:** aggregated API specs update ([94f61c5](https://github.com/G-Core/terraform-provider-gcore/commit/94f61c5302d949b85093cf85baf5b6130f79194c))
* **api:** aggregated API specs update ([df619b9](https://github.com/G-Core/terraform-provider-gcore/commit/df619b9cda3651e889c56989392ba10aac5c5171))
* **api:** aggregated API specs update ([ced4ad4](https://github.com/G-Core/terraform-provider-gcore/commit/ced4ad4058ab7ce52d422d84fba82e09e78694ee))
* **api:** aggregated API specs update ([bfbc3f5](https://github.com/G-Core/terraform-provider-gcore/commit/bfbc3f536d6856b68ce5be364ec2c0faea5f1fa1))
* **api:** aggregated API specs update ([4b4eef7](https://github.com/G-Core/terraform-provider-gcore/commit/4b4eef74f07265febc5037b89fe0fafebbdc0138))
* **api:** aggregated API specs update ([474638d](https://github.com/G-Core/terraform-provider-gcore/commit/474638d2dbb1c9bdcae4f378dc50129abea5ad35))
* **api:** aggregated API specs update ([3669f16](https://github.com/G-Core/terraform-provider-gcore/commit/3669f164d80b00e2742dac2728c5b4cc04583e41))
* **api:** aggregated API specs update ([e999f99](https://github.com/G-Core/terraform-provider-gcore/commit/e999f9995012f912df8791f66287232943148e43))
* **api:** aggregated API specs update ([a93524a](https://github.com/G-Core/terraform-provider-gcore/commit/a93524a72bc3159f1a684769c44a834a2316222d))
* **api:** aggregated API specs update ([83d02fa](https://github.com/G-Core/terraform-provider-gcore/commit/83d02fa40f1b184002461ae8330c556e411c4959))
* **api:** aggregated API specs update ([7b8b82d](https://github.com/G-Core/terraform-provider-gcore/commit/7b8b82d5c4999d166d9e560c96876f8107c2245d))
* **api:** aggregated API specs update ([8098e64](https://github.com/G-Core/terraform-provider-gcore/commit/8098e640753967008f79f820824a4750b28424d4))
* **api:** aggregated API specs update ([114cafb](https://github.com/G-Core/terraform-provider-gcore/commit/114cafb12d76e386e9d91075446c41480637e498))
* **api:** aggregated API specs update ([3faba40](https://github.com/G-Core/terraform-provider-gcore/commit/3faba40a78240faf5d867d369590fdb3cebefc13))
* **api:** aggregated API specs update ([27f553a](https://github.com/G-Core/terraform-provider-gcore/commit/27f553a4ed90c1e4dcb6088e39c18f608d3f87f7))
* **api:** aggregated API specs update ([69be913](https://github.com/G-Core/terraform-provider-gcore/commit/69be9136e1c2327491dbc4280dc5badf661a9231))
* **api:** aggregated API specs update ([35e79cd](https://github.com/G-Core/terraform-provider-gcore/commit/35e79cd35b3a4dff8823850182dae5c873d2aec5))
* **api:** aggregated API specs update ([df14541](https://github.com/G-Core/terraform-provider-gcore/commit/df14541f86c1097cb360794d6cdb4d54699b6332))
* **api:** aggregated API specs update ([73db3a0](https://github.com/G-Core/terraform-provider-gcore/commit/73db3a0ef6c520e0085d9c277e23c7c46594b02c))
* **api:** aggregated API specs update ([6e1e74a](https://github.com/G-Core/terraform-provider-gcore/commit/6e1e74a4a86b6e6668ae3e1924a477c899a69537))
* **api:** aggregated API specs update ([2143572](https://github.com/G-Core/terraform-provider-gcore/commit/2143572a8e49ed7a585453b7f6a61a1d88c17962))
* **api:** aggregated API specs update ([b80cf80](https://github.com/G-Core/terraform-provider-gcore/commit/b80cf80fd98a668117a913216a0efda86dab9ea2))
* **api:** aggregated API specs update ([f4f4c00](https://github.com/G-Core/terraform-provider-gcore/commit/f4f4c00e15745a6aac7779303b0b9cd48640fb54))
* **api:** aggregated API specs update ([72d7ad9](https://github.com/G-Core/terraform-provider-gcore/commit/72d7ad9000b4c6b5b155796536882e05fb8942cc))
* **api:** aggregated API specs update ([b54618c](https://github.com/G-Core/terraform-provider-gcore/commit/b54618cded8e5a45ed275ab949a52d8514bbc70d))
* **api:** aggregated API specs update ([7a34c0b](https://github.com/G-Core/terraform-provider-gcore/commit/7a34c0bb49bd3d3d21b97c6b228b5d3ec6a1159b))
* **api:** aggregated API specs update ([f67dcb2](https://github.com/G-Core/terraform-provider-gcore/commit/f67dcb2e3367069a2ea5f12675197796b4dd860a))
* **api:** aggregated API specs update ([469d17d](https://github.com/G-Core/terraform-provider-gcore/commit/469d17d7535b52d3537bf7af44063343ee2c2e78))
* **api:** aggregated API specs update ([25fe7b6](https://github.com/G-Core/terraform-provider-gcore/commit/25fe7b6177a87979fa718b223c46518e18528f5e))
* **api:** aggregated API specs update ([58bc619](https://github.com/G-Core/terraform-provider-gcore/commit/58bc619bdc7cb16c04608f4c3a73a1290a1df7c4))
* **api:** aggregated API specs update ([d6d3bfe](https://github.com/G-Core/terraform-provider-gcore/commit/d6d3bfede7a129d08b03e1a0f639947578aaabe6))
* **api:** aggregated API specs update ([1629811](https://github.com/G-Core/terraform-provider-gcore/commit/162981136cb73523d383440918c5cb89bbe84456))
* **api:** api update ([800fa27](https://github.com/G-Core/terraform-provider-gcore/commit/800fa27dc1675ad3e1f646910ef562ad77a2d19b))
* **api:** api update ([05d59fd](https://github.com/G-Core/terraform-provider-gcore/commit/05d59fd503c050539c7d579be2f99c2f5396a4db))
* **api:** api update ([c362245](https://github.com/G-Core/terraform-provider-gcore/commit/c362245b76b5d066f7a40e9d2aa278f711c7f407))
* **api:** bump gcore-go SDK dependency ([eeff695](https://github.com/G-Core/terraform-provider-gcore/commit/eeff695c9d1e63f492efc2279737f2540e2de9e7))
* **api:** feat(cdn): enable terraform generation for trusted CA certificates ([#216](https://github.com/G-Core/terraform-provider-gcore/issues/216)) ([f26fc99](https://github.com/G-Core/terraform-provider-gcore/commit/f26fc99627db69aff5f05c7901e970c79aafc031))
* **api:** feat(fastedge): enable TF for apps ([5536489](https://github.com/G-Core/terraform-provider-gcore/commit/5536489c876730bde397f47f317682614bfc7df0))
* **api:** feat(waap): enable terraform generation for waap domains ([#209](https://github.com/G-Core/terraform-provider-gcore/issues/209)) ([a2a0a8b](https://github.com/G-Core/terraform-provider-gcore/commit/a2a0a8b9a9f75acd4d03f4093ff0f37cb952c8c1))
* **api:** manual updates ([6e63441](https://github.com/G-Core/terraform-provider-gcore/commit/6e6344142fee9b75591e68c7234cecce38eec6f9))
* **api:** manual updates ([980633f](https://github.com/G-Core/terraform-provider-gcore/commit/980633fe3ed0390c66b4d9a26d8e80062cf9f7ce))
* **api:** manual updates ([4468400](https://github.com/G-Core/terraform-provider-gcore/commit/44684001219749bb1fc0e402d48f68b66f1a2137))
* **api:** manual updates ([d29d92d](https://github.com/G-Core/terraform-provider-gcore/commit/d29d92d3a711fdcb1e7e55953ae68f7d2f071051))
* **api:** manual updates ([6026123](https://github.com/G-Core/terraform-provider-gcore/commit/6026123200b840bb13f15c7b33c22265954860a2))
* **api:** manual upload of aggregated API specs ([6a638ed](https://github.com/G-Core/terraform-provider-gcore/commit/6a638eddb23494fa3e042cf74b733291e3bb6014))
* **api:** manual upload of aggregated API specs ([7523867](https://github.com/G-Core/terraform-provider-gcore/commit/75238678b63455356c1a719d071c731d81baf6a9))
* **api:** refactor(cdn)!: change type casing from Cdn* to CDN* ([93543ac](https://github.com/G-Core/terraform-provider-gcore/commit/93543acebe5e8f371ec412f88a756b96ea2650d1))
* **api:** revert(cdn): remove client_config subresource ([#207](https://github.com/G-Core/terraform-provider-gcore/issues/207)) ([7adf6a9](https://github.com/G-Core/terraform-provider-gcore/commit/7adf6a9deb517d49db27ecfcf689dc9fb95cc975))
* **api:** trigger build ([8857681](https://github.com/G-Core/terraform-provider-gcore/commit/88576817ab449d8d9b51ea5ef061509ed7d8cdb3))
* **api:** update OpenAPI spec or Stainless config ([77c66aa](https://github.com/G-Core/terraform-provider-gcore/commit/77c66aa29a43bf37e7ceeeada30a3c7bfda8b085))
* **api:** update via SDK Studio ([029e669](https://github.com/G-Core/terraform-provider-gcore/commit/029e669973f3051211b02ceae67433568ff8de45))
* **api:** update via SDK Studio ([af5558f](https://github.com/G-Core/terraform-provider-gcore/commit/af5558ffe028767ae0c9eb8a4c0088b42692bf8e))
* **cdn:** add client_config subresource for terraform ([3d3f959](https://github.com/G-Core/terraform-provider-gcore/commit/3d3f95964c6e04efdcaced716078f456f8786988))
* **cdn:** add custom behaviors for cdn_certificate resource ([#105](https://github.com/G-Core/terraform-provider-gcore/issues/105)) ([fc7346b](https://github.com/G-Core/terraform-provider-gcore/commit/fc7346bcaea8e708f9a75035852a8e3d35d1254e))
* **cdn:** enable terraform for cdn_origingroup ([383474b](https://github.com/G-Core/terraform-provider-gcore/commit/383474bd651044f48effc12401577b397637d0a5))
* **cdn:** enable terraform generation for cdn rules ([9bccb16](https://github.com/G-Core/terraform-provider-gcore/commit/9bccb160b3e6909116662e15104017fdf468ef04))
* **cdn:** enable terraform generation for SSL certificates ([398e4e0](https://github.com/G-Core/terraform-provider-gcore/commit/398e4e09153ea7a307c1f25cd5bda23fa7828230))
* **cdn:** enable terraform generation for trusted CA certificates ([9aaeff5](https://github.com/G-Core/terraform-provider-gcore/commit/9aaeff53161177dc8e9c3b16aed9456f4b84dd44))
* **cdn:** override name for cdn_resource_rule to avoid duplication ([d9260dd](https://github.com/G-Core/terraform-provider-gcore/commit/d9260dd3ebb30733776c613225fcde1d9f348f22))
* **cdn:** override name for cdn_resource_rule to avoid duplication ([#224](https://github.com/G-Core/terraform-provider-gcore/issues/224)) ([5782b27](https://github.com/G-Core/terraform-provider-gcore/commit/5782b27c550205645d5f40022b1202d4c948b13c))
* **cloud_security_group_rule:** implement CRUD with v2 async endpoints ([#131](https://github.com/G-Core/terraform-provider-gcore/issues/131)) ([028bf41](https://github.com/G-Core/terraform-provider-gcore/commit/028bf41097febd56794aa41e4713242a7af2a3c6))
* **cloud:** add custom code for TF file shares and access rules ([b661e27](https://github.com/G-Core/terraform-provider-gcore/commit/b661e27e275cc5da1dbc9f1b5b12ac79170f4aad))
* **cloud:** add full lifecycle management for Terraform k8s clusters ([ce2ba0f](https://github.com/G-Core/terraform-provider-gcore/commit/ce2ba0fe7e1d2aa1123c0c7703c8f7dbe3c0da0a))
* **cloud:** add terraform resource and data source for rrsets ([a7ad415](https://github.com/G-Core/terraform-provider-gcore/commit/a7ad4150ccaaec80e7fc779d4845c7b2fb05395f))
* **cloud:** add terraform support for GPU virtual images ([0246a11](https://github.com/G-Core/terraform-provider-gcore/commit/0246a1121ce04ac7d48742cf0e31c79a23172b5e))
* **cloud:** add terraform support for inference deployment ([9ecd9e5](https://github.com/G-Core/terraform-provider-gcore/commit/9ecd9e5eb5283f40a6473b1543ade9d1a6d0eda0))
* **cloud:** convert cloud_secret payload fields to write-only attributes ([de0aa6e](https://github.com/G-Core/terraform-provider-gcore/commit/de0aa6e052759d7fb090e72d6af8e51eae2e7ae3))
* **cloud:** enable terraform for security group rules (v2) ([6c3be87](https://github.com/G-Core/terraform-provider-gcore/commit/6c3be87704bbcdc3c93da43c2cddca652784e748))
* **cloud:** enable terraform resource for gcore_image ([a8a4a82](https://github.com/G-Core/terraform-provider-gcore/commit/a8a4a82b6135ec22072bc25e494d2a94d82eb703))
* **cloud:** enable TF for file share and file share access rules ([22594f5](https://github.com/G-Core/terraform-provider-gcore/commit/22594f5255eb6c103e815cfcf93b17c55c01ffeb))
* **cloud:** enable tf for GPU virtual clusters ([fd0ea9b](https://github.com/G-Core/terraform-provider-gcore/commit/fd0ea9b896243d2f8f3848819340cbe967a73ca6))
* **cloud:** enable TF for inference registry credentials ([b16da92](https://github.com/G-Core/terraform-provider-gcore/commit/b16da92f375efbc20d611d85ce6cdb7a1f8a864d))
* **cloud:** implement custom gpu virtual cluster resource with drift prevention and import support ([d73b74c](https://github.com/G-Core/terraform-provider-gcore/commit/d73b74ce12048aaf3b9adad83f77786b4c503d04))
* **cloud:** make cloud_instance password a write-only field ([19d5952](https://github.com/G-Core/terraform-provider-gcore/commit/19d5952ac4586d9fd5b3986ce3d079dd64431805))
* **cloud:** make inference secret data fields write-only ([2137e25](https://github.com/G-Core/terraform-provider-gcore/commit/2137e2562581fd8141db80d4ac475951c3712306))
* **cloud:** make password a write-only field ([20af6b5](https://github.com/G-Core/terraform-provider-gcore/commit/20af6b5b3801d4732a0206452520c32ee1806ef9))
* **cloud:** migrate instance resource to use *AndPoll methods ([#35](https://github.com/G-Core/terraform-provider-gcore/issues/35)) ([b40c713](https://github.com/G-Core/terraform-provider-gcore/commit/b40c71331f5810622da49350a3de4672c15a38dd))
* **cloud:** remove ddos_profile from k8s cluster resource and data source ([81e11c8](https://github.com/G-Core/terraform-provider-gcore/commit/81e11c836e381ec3f45096dd0ac68a66984ff567))
* **cloud:** set k8s TF id_property as name ([a46659a](https://github.com/G-Core/terraform-provider-gcore/commit/a46659af0af200f6e66c2bfc392ded2ca897b765))
* **cloud:** use *AndPoll methods for GPU bare metal images ([aeeb7c0](https://github.com/G-Core/terraform-provider-gcore/commit/aeeb7c020d1b6f2093085a0993dc9575d7d955d6))
* **cloud:** use create and update v2 endpoints for security groups ([c941448](https://github.com/G-Core/terraform-provider-gcore/commit/c94144810cd948c6a4a90a7aafbfc7e81c1cb45d))
* **cloud:** use data_wo_version to trigger updates ([030aed0](https://github.com/G-Core/terraform-provider-gcore/commit/030aed0dbd0bd15f7aebc977d36c7e1eb35813fb))
* **cloud:** use Replace method in inference secrets TF updates ([9a34c7e](https://github.com/G-Core/terraform-provider-gcore/commit/9a34c7e7f943a95c5e09d256d4b2a3bd226e49d1))
* **cloud:** use Replace method in registry credentials TF updates ([63ee17a](https://github.com/G-Core/terraform-provider-gcore/commit/63ee17a88087b33e49264ba515b88cb0eea02965))
* **cloud:** use UpdateAndPoll for floating IP updates ([8c38c8a](https://github.com/G-Core/terraform-provider-gcore/commit/8c38c8a55c793a598c9ca76dc83863efbd2f25df))
* **cloud:** use v2 endpoint for floating IPs updates ([56245ad](https://github.com/G-Core/terraform-provider-gcore/commit/56245adfbe5d42e9ca8a8bbb5021aca12b48ed38))
* **dns_zone:** MetaStringType for meta fields, DNSSEC, import ([#114](https://github.com/G-Core/terraform-provider-gcore/issues/114)) ([8df76ab](https://github.com/G-Core/terraform-provider-gcore/commit/8df76ab7923aa4221543b2dfa97fe4d6305c2e06))
* **dns:** add gcore_dns_zone_rrset resource with full CRUD and import support ([#60](https://github.com/G-Core/terraform-provider-gcore/issues/60)) ([9fbf01c](https://github.com/G-Core/terraform-provider-gcore/commit/9fbf01c1fe21851c40f1ded52bff71cf7aef2b0d))
* **dns:** enable terraform code generation for gcore_dns_network_mapping ([8b8dbb1](https://github.com/G-Core/terraform-provider-gcore/commit/8b8dbb12a44140995bd13e64a51abddcf8555951))
* **dns:** enable terraform code generation for gcore_dns_zone ([c382d75](https://github.com/G-Core/terraform-provider-gcore/commit/c382d75d8472d505607d3d8b84ddd817193c0a50))
* **dns:** MetaStringType, update support, import for dns_zone_rrset ([f70568d](https://github.com/G-Core/terraform-provider-gcore/commit/f70568dd962a1451a69f0cb382127e575865b403))
* **dns:** MetaStringType, update support, import for dns_zone_rrset ([#115](https://github.com/G-Core/terraform-provider-gcore/issues/115)) ([a42cea8](https://github.com/G-Core/terraform-provider-gcore/commit/a42cea8412620c79a10691f29bcd97429af83c50))
* **dns:** register DNS data sources in provider ([0dc4d70](https://github.com/G-Core/terraform-provider-gcore/commit/0dc4d70f3f606ed3fa9da42fbf495b4a131c298c))
* **docs:** migrate examples and templates from old provider ([c873c86](https://github.com/G-Core/terraform-provider-gcore/commit/c873c862a145b79035f98d2cbe863cf3e919d9cf))
* **fastedge:** enable fastedge_secret resource generation ([75993a2](https://github.com/G-Core/terraform-provider-gcore/commit/75993a2a7129b0e092d1b6472931b1634f534abf))
* **fastedge:** enable terraform for binaries resource ([9617307](https://github.com/G-Core/terraform-provider-gcore/commit/9617307c1e35f9e49514040c4b97229eca91af6e))
* **fastedge:** enable TF for apps ([b072f0e](https://github.com/G-Core/terraform-provider-gcore/commit/b072f0ed134119ea35792051c8726a08812e8eeb))
* **terraform:** enable fastedge_secret resource generation ([#202](https://github.com/G-Core/terraform-provider-gcore/issues/202)) ([a589e87](https://github.com/G-Core/terraform-provider-gcore/commit/a589e8772447ddeffdbb93dfb59514a4d6f3413a))
* **waap:** enable terraform support for WAAP resources ([11ccd46](https://github.com/G-Core/terraform-provider-gcore/commit/11ccd46adac66ba80008525dfe9f15db682d496c))


### Bug Fixes

* **api:** revert unnecessary changes to additionalProperties ([6ad95d3](https://github.com/G-Core/terraform-provider-gcore/commit/6ad95d30737eda0095124ba6ed346f3ed672d2ce))
* bump gcore-go to v0.26.0 for GPUVirtual support ([084325d](https://github.com/G-Core/terraform-provider-gcore/commit/084325d7e572ec9955527ca049035a0808a2293d))
* **cdn:** fix drift detection for computed fields ([#70](https://github.com/G-Core/terraform-provider-gcore/issues/70)) ([0ddad04](https://github.com/G-Core/terraform-provider-gcore/commit/0ddad04dd04b408b5b7d7043077483cffa4186c1))
* **cdn:** fix drift detection for computed fields ([#82](https://github.com/G-Core/terraform-provider-gcore/issues/82)) ([e7a687c](https://github.com/G-Core/terraform-provider-gcore/commit/e7a687c1be6d118eb5ff2eae50282eb26bc1868b))
* **cdn:** resolve CDN resource drift detection and update serialization issues ([7d899ff](https://github.com/G-Core/terraform-provider-gcore/commit/7d899ff5b759ed9602cb97c8d1631a3fa2543eea))
* **cdn:** update new model names in origin groups tf resource ([dc4eb9b](https://github.com/G-Core/terraform-provider-gcore/commit/dc4eb9bf0f8ee8d9a22337355fecb1927d732167))
* change repo for fossa to G-Core/terraform-provider-gcore ([08fb46a](https://github.com/G-Core/terraform-provider-gcore/commit/08fb46aa7d1e3a38429c380cef72d384779a3e94))
* **client:** mark count in offset pagination as required ([e81c301](https://github.com/G-Core/terraform-provider-gcore/commit/e81c301c63a87ddd33b11a2cdc4585accfdeb1a2))
* **cloud:** add missing json tag to access rule ids ([5611e36](https://github.com/G-Core/terraform-provider-gcore/commit/5611e369bf5c0ba48e89e5ae0bd0f51fa4bfa4ff))
* **cloud:** instance image, use UploadAndPoll and DeleteAndPoll for async operations ([#74](https://github.com/G-Core/terraform-provider-gcore/issues/74)) ([950913f](https://github.com/G-Core/terraform-provider-gcore/commit/950913f21b8ec922ab47532babf4e742a569faa8))
* **cloud:** instance, use UpdateAndPoll for floating IP assign/unassign ([c6e0d20](https://github.com/G-Core/terraform-provider-gcore/commit/c6e0d203194a357606f26c6ff52071a71a3678b8))
* **cloud:** k8s references from k8 to k8s ([d28619b](https://github.com/G-Core/terraform-provider-gcore/commit/d28619b0cbf319c58067eb2e95176d6c615851eb))
* **cloud:** make file share name required in schema ([377c2cd](https://github.com/G-Core/terraform-provider-gcore/commit/377c2cda53375e95d6ac3fffc05bb8eaa350f720))
* **cloud:** remove listeners from LB resource and deprecated timeout fields ([#163](https://github.com/G-Core/terraform-provider-gcore/issues/163)) ([ff40138](https://github.com/G-Core/terraform-provider-gcore/commit/ff40138579b39586a74af69ecb842ea9714dd0ce))
* **cloud:** rename load balancer pool member methods to create/delete ([89f94d0](https://github.com/G-Core/terraform-provider-gcore/commit/89f94d0a3b78391a8f893ce5087f59288ecd83bd))
* **cloud:** replace PUT /cloud/v1/l7policies with PATCH ([6754672](https://github.com/G-Core/terraform-provider-gcore/commit/6754672ba84104fea664167f7f6aa9b522743412))
* **cloud:** restructure to be gpu_virtual.clusters ([7156ede](https://github.com/G-Core/terraform-provider-gcore/commit/7156ede01104a8f6fdc45d969ff78bbae0079c2b))
* **cloud:** use load balancer pool member methods new names ([cc2e1b0](https://github.com/G-Core/terraform-provider-gcore/commit/cc2e1b0e6f6be902054dacf7f0f8a6f3ccd4aa9f))
* compilation of Go SDK call with binary file upload ([300124f](https://github.com/G-Core/terraform-provider-gcore/commit/300124f73fc9cc44d0ff079cd5eb4811dbf77817))
* correct partial route deletion condition and preserve planned routes ([0583bc2](https://github.com/G-Core/terraform-provider-gcore/commit/0583bc254f27cd1e3a1d45bb1fcb41843a2faa96))
* correctly mark a subset of fields shared between create and update calls as required ([1a33b65](https://github.com/G-Core/terraform-provider-gcore/commit/1a33b65ab19dd29cf7ec516a734396d6019506ad))
* datasource timeout schema ([3ff0839](https://github.com/G-Core/terraform-provider-gcore/commit/3ff0839861a71c13e440d4ed2911f06d8600ccbb))
* dedup data source query attributes ([8b16e26](https://github.com/G-Core/terraform-provider-gcore/commit/8b16e262db6b038a63dcadc9cdd686d54103ce45))
* **dns:** add id_property to dns_zone for correct identifier mapping ([2903803](https://github.com/G-Core/terraform-provider-gcore/commit/29038037c0ece104ce0a37221fefb1c22f59725b))
* **dns:** move id_path_param to method level for dns_zone ([c450be8](https://github.com/G-Core/terraform-provider-gcore/commit/c450be8146f42b3b0b5503356261d1b04370d1ad))
* ensure derived request attribute schemas conform to the upstream configurability overrides ([ad847fa](https://github.com/G-Core/terraform-provider-gcore/commit/ad847facfa3bc63933480dddf136341c6cb025bf))
* **fastedge:** add sensitive field handling and force delete ([#94](https://github.com/G-Core/terraform-provider-gcore/issues/94)) ([a98e71e](https://github.com/G-Core/terraform-provider-gcore/commit/a98e71e399c2d19e413b54d07d178d51238d9e47))
* **fastedge:** implement file upload for gcore_fastedge_binary resource ([#84](https://github.com/G-Core/terraform-provider-gcore/issues/84)) ([3ec760c](https://github.com/G-Core/terraform-provider-gcore/commit/3ec760c65c13dc58efb2462bf03fcd6ed9d8daf1))
* **fastedge:** prevent state drift, validate either template or binary, use PUT instead of PATCH for apps ([ffb7e8e](https://github.com/G-Core/terraform-provider-gcore/commit/ffb7e8e7723c85729f3d0beb453abc1e52b9b6e4))
* **fastedge:** remove linting error with temporary workaround ([425c36f](https://github.com/G-Core/terraform-provider-gcore/commit/425c36f207d0a70141cbbb34b110555fc4a52427))
* **fastedge:** remove readOnly name from app_store required fields ([58e03a7](https://github.com/G-Core/terraform-provider-gcore/commit/58e03a7827ca9ca36d8e3ae577f36a47ea5f46fa))
* improve acctest script package handling ([20fda80](https://github.com/G-Core/terraform-provider-gcore/commit/20fda80ea20d8dd3fcfe45a54b61616e1c68f16f))
* list style data sources should always have id value populated ([c8788d1](https://github.com/G-Core/terraform-provider-gcore/commit/c8788d1d6b52b93bbdedf0295c462bb8dac2e18d))
* prefer named identifier field over id alias, missing ImportStates in certain resources ([b7dfae1](https://github.com/G-Core/terraform-provider-gcore/commit/b7dfae1b5b0d13f5a2df90df0c6c625a87b88043))
* **provider:** remove duplicate cdn_origin_group data source registration ([3148739](https://github.com/G-Core/terraform-provider-gcore/commit/3148739a3bbae08233136bc4fec1f9d98bd8b1b8))
* remove task_id, creator_task_id, and tasks fields from schemas ([9005503](https://github.com/G-Core/terraform-provider-gcore/commit/90055032a2d783423429d50753dbe54dea6bbc92))
* replace stale stainless-sdks import paths with G-Core module path ([247bdae](https://github.com/G-Core/terraform-provider-gcore/commit/247bdae97c2f3853dd20be61cf7540cf2fd2ed91))
* resolve compilation errors from tags jsontypes migration ([ebabfac](https://github.com/G-Core/terraform-provider-gcore/commit/ebabfac3f7f98b1e26a335c09c03751d10819d7f))
* resolve configuration drift in load balancer and listener resources ([1aa974c](https://github.com/G-Core/terraform-provider-gcore/commit/1aa974c24afcc04fd295c1f8d3966f3e63304dd2))
* revert tags jsontypes migration fixes ([ebc3bbc](https://github.com/G-Core/terraform-provider-gcore/commit/ebc3bbc373c505845d9aa83695343c4564737c5c))
* **router:** preserve interfaces field after PATCH response unmarshal ([d3cf4be](https://github.com/G-Core/terraform-provider-gcore/commit/d3cf4be84c9cc17aac29e34bc8655143eae2eeda))
* spurious update plans for float attributes after import ([94fc8c9](https://github.com/G-Core/terraform-provider-gcore/commit/94fc8c92dc7e897915aafda72f12c64f06330ee7))
* **terraform:** add missing imports in provider.go ([598012f](https://github.com/G-Core/terraform-provider-gcore/commit/598012f75adf8c6c4976fb40fdcf97f0eca30ee4))
* **terraform:** remove duplicate resource references in provider ([56f8003](https://github.com/G-Core/terraform-provider-gcore/commit/56f80034d46eb544ec52fbca7849f8274b96eb0a))
* **terraform:** remove more duplicate resource references in provider ([e7a3616](https://github.com/G-Core/terraform-provider-gcore/commit/e7a36166753675f9cfd60e821cb056ebb3834bdf))
* **types:** correctly define false enum ([93f927f](https://github.com/G-Core/terraform-provider-gcore/commit/93f927f71399fb77cf12aa2cb089ff54d1c2091d))
* update stale module path imports for next branch ([2c8f2d1](https://github.com/G-Core/terraform-provider-gcore/commit/2c8f2d1f375fe3cf22a7175ee093662888ffcebf))
* update stale module path imports for next branch ([1f78c8c](https://github.com/G-Core/terraform-provider-gcore/commit/1f78c8cf48e6985ce2e2ff2886d018ff3d125a2c))


### Chores

* add .env and .envrc to gitignore ([ef53fdc](https://github.com/G-Core/terraform-provider-gcore/commit/ef53fdc49db41ed6230988007e347a7069ecbfa7))
* bump dependency version ([ac1c8db](https://github.com/G-Core/terraform-provider-gcore/commit/ac1c8db42d757f9bcd4dfdc289b9c45284ab7a9d))
* **cdn:** fix name of gcore_cdn_resource_rule ([bb1db5c](https://github.com/G-Core/terraform-provider-gcore/commit/bb1db5cfbff60dfb129e8eabcfab24c2093f7ac3))
* change naming for POST, PUT, PATCH, DELETE models ([9e5a935](https://github.com/G-Core/terraform-provider-gcore/commit/9e5a93570ec9b060b62a1e446503f699f78b33ac))
* **cloud:** disable terraform for inference deployment ([cea1c33](https://github.com/G-Core/terraform-provider-gcore/commit/cea1c33c5d821b844930b690f5619236663ee0cb))
* **docs:** add missing descriptions ([3390674](https://github.com/G-Core/terraform-provider-gcore/commit/3390674fe1cf93c0aa8c46c1b5654347aec60506))
* **docs:** regenerate docs to sync with current schemas ([#165](https://github.com/G-Core/terraform-provider-gcore/issues/165)) ([46ba2ce](https://github.com/G-Core/terraform-provider-gcore/commit/46ba2cef22f7c578e55e7022c90855dc23203c90))
* ensure tests build as part of lint step ([3cde506](https://github.com/G-Core/terraform-provider-gcore/commit/3cde506433ae0634d2ca6973600ef8a8f092c1f5))
* fix conflicts ([3e705f2](https://github.com/G-Core/terraform-provider-gcore/commit/3e705f29e2a42440d9cb9a4cf57b53bc50572e3f))
* **internal:** codegen related update ([2cf5a3f](https://github.com/G-Core/terraform-provider-gcore/commit/2cf5a3f4065906b910d7d8017a8fc8bc158fe24b))
* **internal:** codegen related update ([bab0f49](https://github.com/G-Core/terraform-provider-gcore/commit/bab0f4935de3acebf9f2ed964e5e016183cff99e))
* **internal:** codegen related update ([804e8e5](https://github.com/G-Core/terraform-provider-gcore/commit/804e8e5a17a3dee3cbf00382029ef195d8bc1565))
* **internal:** codegen related update ([5f67f0a](https://github.com/G-Core/terraform-provider-gcore/commit/5f67f0a40e60b12a11aa60cdd93de46729119e65))
* **internal:** codegen related update ([6310816](https://github.com/G-Core/terraform-provider-gcore/commit/63108163a714088d7c3b32f646f874038a11ea33))
* **internal:** codegen related update ([01f72f1](https://github.com/G-Core/terraform-provider-gcore/commit/01f72f14956e11d8b48f4ed8a69ba13e15348931))
* **internal:** codegen related update ([75cff8d](https://github.com/G-Core/terraform-provider-gcore/commit/75cff8de6d30438012526e104267e3bb942b8a07))
* **internal:** codegen related update ([daac42f](https://github.com/G-Core/terraform-provider-gcore/commit/daac42f8b23472694f5dc06fe01265a219f2a76a))
* **internal:** codegen related update ([93b5016](https://github.com/G-Core/terraform-provider-gcore/commit/93b501654609fc1b84d86930cb0a4cfea1498abb))
* **internal:** codegen related update ([8379f85](https://github.com/G-Core/terraform-provider-gcore/commit/8379f859c08273e9ef1c7b1d0a7f831acf44cfde))
* **internal:** codegen related update ([8fedffd](https://github.com/G-Core/terraform-provider-gcore/commit/8fedffda65cccaca9021355cb458e9a29abe3bbf))
* **internal:** codegen related update ([6f75ef2](https://github.com/G-Core/terraform-provider-gcore/commit/6f75ef286fd44fb154b6a57500391e4f40e45161))
* **internal:** codegen related update ([40cd9c5](https://github.com/G-Core/terraform-provider-gcore/commit/40cd9c53517d289c9cdf3c75b535e66c95ac486e))
* **internal:** codegen related update ([db0877c](https://github.com/G-Core/terraform-provider-gcore/commit/db0877cb812c96b75bc39bb979528c1e884b5720))
* **internal:** codegen related update ([d5373a9](https://github.com/G-Core/terraform-provider-gcore/commit/d5373a9eba8622eae97b6884ae07214b5f40cf24))
* **internal:** codegen related update ([f1068a5](https://github.com/G-Core/terraform-provider-gcore/commit/f1068a508af911107cc6fb7e8e54ba2830d8ef39))
* **internal:** update `actions/checkout` version ([7add8eb](https://github.com/G-Core/terraform-provider-gcore/commit/7add8ebdfa926c0d99add0030627e9c4e78e6d3c))
* **test:** do not count install time for mock server timeout ([c22bbda](https://github.com/G-Core/terraform-provider-gcore/commit/c22bbda93c5668fb48024ea442fa9a54440ae5bc))
* update SDK settings ([6c43c61](https://github.com/G-Core/terraform-provider-gcore/commit/6c43c61caec2fa5260d10aab600d7a7bb9eda878))


### Documentation

* add descriptions to client settings options ([72c85ad](https://github.com/G-Core/terraform-provider-gcore/commit/72c85adf79de7619da54a60f4df4fc13483bc15f))
* add schema-behaviors skill ([92e8e03](https://github.com/G-Core/terraform-provider-gcore/commit/92e8e0383b1150020f2d7ad3c6179fc835a3f981))
* bump minimum Terraform CLI version to 1.11 ([6beba87](https://github.com/G-Core/terraform-provider-gcore/commit/6beba87795af167ff2121ebab23fd14de824f430))
* **cloud:** update cloud_instance and cloud_secret with latest schemas ([7984796](https://github.com/G-Core/terraform-provider-gcore/commit/79847966ab7ff3e5dbe3a4d728b8352b4ed4e86c))
* replace stale stainless-sdks import paths ([44e3d90](https://github.com/G-Core/terraform-provider-gcore/commit/44e3d9041b17f4f11e4a37664f7be3ed802cd0cb))
* streamline text in index.md ([31bb803](https://github.com/G-Core/terraform-provider-gcore/commit/31bb803a2f9f0cdf3ceb42fbcd9fd4a2e2ee3825))


### Refactors

* move generic plan modifiers to shared package ([53fcc2d](https://github.com/G-Core/terraform-provider-gcore/commit/53fcc2df1804b4470c683415f9d4cd41632c24b9))
* remove unnecessary WithJSONSet and clean up comments ([fdccfc6](https://github.com/G-Core/terraform-provider-gcore/commit/fdccfc653f82286251b3cb8f9c8fd7070c0b775f))
* **router:** optimize Update flow with consolidated PATCH and correct API call ordering ([2d221cd](https://github.com/G-Core/terraform-provider-gcore/commit/2d221cd5bab61e6136e7202267d4882b541f33cb))
* **router:** remove unnecessary preservedInterfaces logic ([6bf5261](https://github.com/G-Core/terraform-provider-gcore/commit/6bf5261985326a24f2b66a62899178be9eee35e9))
* **router:** use ([8ea1064](https://github.com/G-Core/terraform-provider-gcore/commit/8ea10645ad44f1b7d6fcbf77f0756adbf02d5a88))
