package licensing

import (
	"regexp"
	"strings"
)

// All mapping keys must be uppercase
var mapping = map[string]string{
	// Simple mappings (i.e. that could be parsed by SpdxExpression.parse, at least without space)
	// modified from https://github.com/oss-review-toolkit/ort/blob/fc5389c2cfd9c8b009794c8a11f5c91321b7a730/utils/spdx/src/main/resources/simple-license-mapping.yml

	// Ambiguous simple mappings (mapping reason not obvious without additional information)
	"AFL":          AFL30,
	"AGPL":         AGPL30,
	"ALV2":         Apache20,
	"APACHE":       Apache20,
	"APACHE-STYLE": Apache20,
	"ARTISTIC":     Artistic20,
	"ASL":          Apache20,
	"BSD":          BSD3Clause,
	"BSD*":         BSD3Clause,
	"BSD-LIKE":     BSD3Clause,
	"BSD-STYLE":    BSD3Clause,
	"CDDL":         CDDL10,
	"ECLIPSE":      EPL10,
	"EPL":          EPL10,
	"EUPL":         EUPL10,
	"FDL":          GFDL13,
	"GFDL":         GFDL13,
	"GPL":          GPL20,
	"LGPL":         LGPL20,
	"MPL":          MPL20,
	"NETSCAPE":     NPL11,
	"PYTHON":       Python20,

	// Non-ambiguous simple mappings
	"AFL-2":                       AFL20,
	"AFL-2.1":                     AFL21,
	"AFL2":                        AFL20,
	"AFL2.0":                      AFL20,
	"AFL2.1":                      AFL21,
	"APACHE-2":                    Apache20,
	"APACHE-2.0":                  Apache20,
	"APACHE2":                     Apache20,
	"APACHE2.0":                   Apache20,
	"APACHEV2":                    Apache20,
	"APL2":                        Apache20,
	"APLV2.0":                     Apache20,
	"BOOST":                       BSL10,
	"BOUNCY":                      MIT,
	"BOUNCY-LICENSE":              MIT,
	"BSD-2-CLAUSE":                BSD2Clause,
	"BSD-3":                       BSD3Clause,
	"BSD-3-CLAUSE":                BSD3Clause,
	"BSD-4-CLAUSE":                BSD4Clause,
	"BSD2":                        BSD2Clause,
	"BSD3":                        BSD3Clause,
	"BSL":                         BSL10,
	"BSL1.0":                      BSL10,
	"CC0":                         CC010,
	"CDDL1.0":                     CDDL10,
	"CDDL1.1":                     CDDL11,
	"CPAL":                        CPAL10,
	"CPL":                         CPL10,
	"EDL-1.0":                     BSD3Clause,
	"EPL1.0":                      EPL10,
	"EPL2.0":                      EPL20,
	"EUPL1.0":                     EUPL10,
	"EUPL1.1":                     EUPL11,
	"EXPAT":                       MIT,
	"GO":                          BSD3Clause,
	"GPL-1":                       GPL10,
	"GPL-1.0":                     GPL10,
	"GPL-2":                       GPL20,
	"GPL-2+-WITH-BISON-EXCEPTION": GPL20withbisonexception,
	"GPL-2.0":                     GPL20,
	"GPL-3":                       GPL30,
	"GPL-3+-WITH-BISON-EXCEPTION": GPL20withbisonexception,
	"GPLV2+CE":                    GPL20withclasspathexception,
	"GPL-3.0":                     GPL30,
	"GPL1":                        GPL10,
	"GPL2":                        GPL20,
	"GPL3":                        GPL30,
	"GPLV1":                       GPL10,
	"GPLV2":                       GPL20,
	"GPLV3":                       GPL30,
	"HSQLDB":                      BSD3Clause,
	"ISCL":                        ISC,
	"JQUERY":                      MIT,
	"LGPL-2":                      LGPL20,
	"LGPL-2.0":                    LGPL20,
	"LGPL-2.1":                    LGPL21,
	"LGPL-3":                      LGPL30,
	"LGPL-3.0":                    LGPL30,
	"LGPL2":                       LGPL21,
	"LGPL3":                       LGPL30,
	"LGPLV2":                      LGPL21,
	"LGPLV2.1":                    LGPL21,
	"LGPLV3":                      LGPL30,
	// MIT No Attribution (MIT-0) is not yet supported by google/licenseclassifier
	"MIT-0":      MIT,
	"MIT-LIKE":   MIT,
	"MIT-STYLE":  MIT,
	"MPL-2":      MPL20,
	"MPL1":       MPL10,
	"MPL1.0":     MPL10,
	"MPL2":       MPL20,
	"MPL2.0":     MPL20,
	"MPLV2":      MPL20,
	"MPLV2.0":    MPL20,
	"POSTGRESQL": PostgreSQL,
	"RUBY":       Ruby,
	"UNLICENSE":  Unlicense,
	"UNLICENSED": Unlicense,
	"W3CL":       W3C,
	"WTF":        WTFPL,
	"ZLIB":       Zlib,
	"ZOPE":       ZPL21,

	// Non simple declared mappings
	// modified from https://github.com/oss-review-toolkit/ort/blob/fc5389c2cfd9c8b009794c8a11f5c91321b7a730/utils/spdx/src/main/resources/declared-license-mapping.yml

	// Ambiguous declared mappings (mapping reason not obvious without additional information)
	"ACADEMIC FREE LICENSE (AFL)":                         AFL21,
	"APACHE SOFTWARE":                                     Apache20,
	"APACHE SOFTWARE LICENSES":                            Apache20,
	"APPLE PUBLIC SOURCE":                                 APSL10,
	"BSD 4-CLAUSE":                                        BSD3Clause,
	"BSD SOFTWARE":                                        BSD2Clause,
	"BSD STYLE":                                           BSD3Clause,
	"COMMON DEVELOPMENT AND DISTRIBUTION":                 CDDL10,
	"CREATIVE COMMONS - BY":                               CCBY30,
	"CREATIVE COMMONS ATTRIBUTION":                        CCBY30,
	"CREATIVE COMMONS":                                    CCBY30,
	"ECLIPSE PUBLIC LICENSE (EPL)":                        EPL10,
	"GENERAL PUBLIC LICENSE (GPL)":                        GPL20,
	"GNU FREE DOCUMENTATION LICENSE (FDL)":                GFDL13,
	"GNU GENERAL PUBLIC LICENSE (GPL)":                    GPL30,
	"GNU GPL":                                             GPL20,
	"GNU LESSER":                                          LGPL21,
	"GNU LESSER GENERAL PUBLIC":                           LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE (LGPL)":            LGPL21,
	"GNU LESSER PUBLIC":                                   LGPL21,
	"GNU LGPL":                                            LGPL21,
	"GNU LIBRARY OR LESSER GENERAL PUBLIC LICENSE (LGPL)": LGPL21,
	"GNU PUBLIC":                                          GPL20,
	"GPL (WITH DUAL LICENSING OPTION)":                    GPL20,
	"GPLV2 WITH EXCEPTIONS":                               GPL20withclasspathexception,
	"INDIVIDUAL BSD":                                      BSD3Clause,
	"LESSER GENERAL PUBLIC LICENSE (LGPL)":                LGPL21,
	"LGPL WITH EXCEPTIONS":                                LGPL30,
	"LGPLV2 OR LATER":                                     LGPL21,
	"LPGL, SEE LICENSE FILE.":                             LGPL30,
	"MOZILLA PUBLIC":                                      MPL20,

	// Non-ambiguous declared mappings
	"(NEW) BSD":                       BSD3Clause,
	"2-CLAUSE BSD":                    BSD2Clause,
	"2-CLAUSE BSDL":                   BSD2Clause,
	"3-CLAUSE BDSL":                   BSD3Clause,
	"3-CLAUSE BSD":                    BSD3Clause,
	"ACADEMIC FREE LICENSE (AFL-2.1)": AFL21,
	"AFFERO GENERAL PUBLIC LICENSE (AGPL) V. 3": AGPL30,
	"AGPL V3":                          AGPL30,
	"AL 2.0":                           Apache20,
	"APACHE VERSION 2.0, JANUARY 2004": Apache20,
	"APACHE 2 STYLE":                   Apache20,
	"APACHE 2":                         Apache20,
	"APACHE 2.0":                       Apache20,
	"APACHE LICENSE (2.0)":             Apache20,
	"APACHE LICENSE (V2.0)":            Apache20,
	"APACHE LICENSE 2":                 Apache20,
	"APACHE LICENSE 2.0":               Apache20,
	"APACHE LICENSE V2":                Apache20,
	"APACHE LICENSE V2.0":              Apache20,
	"APACHE LICENSE VERSION 2":         Apache20,
	"APACHE LICENSE VERSION 2.0":       Apache20,
	"APACHE LICENSE, 2.0":              Apache20,
	"APACHE LICENSE, ASL VERSION 2.0":  Apache20,
	"APACHE LICENSE, V2 OR LATER":      Apache20,
	"APACHE LICENSE, V2.0 OR LATER":    Apache20,
	"APACHE LICENSE, VERSION 2":        Apache20,
	"APACHE LICENSE, VERSION 2.0 (HTTP://WWW.APACHE.ORG/LICENSES/LICENSE-2.0)": Apache20,
	"APACHE LICENSE, VERSION 2.0":                        Apache20,
	"APACHE LICENSE,VERSION 2.0":                         Apache20,
	"APACHE PUBLIC LICENSE 2.0":                          Apache20,
	"APACHE SOFTWARE LICENSE (APACHE-2.0)":               Apache20,
	"APACHE SOFTWARE LICENSE - VERSION 2.0":              Apache20,
	"APACHE SOFTWARE LICENSE 2.0":                        Apache20,
	"APACHE SOFTWARE LICENSE, VERSION 1.1":               Apache11,
	"APACHE SOFTWARE LICENSE, VERSION 2":                 Apache20,
	"APACHE SOFTWARE LICENSE, VERSION 2.0":               Apache20,
	"APACHE V 2.0":                                       Apache20,
	"APACHE V2":                                          Apache20,
	"APACHE V2.0":                                        Apache20,
	"APACHE VERSION 2.0":                                 Apache20,
	"APACHE, VERSION 2.0":                                Apache20,
	"APACHE-2.0 */ &#39; &QUOT; &#X3D;END --":            Apache20,
	"ARTISTIC 2.0":                                       Artistic20,
	"ARTISTIC LICENSE V2.0":                              Artistic20,
	"ASF 2.0":                                            Apache20,
	"ASL 2":                                              Apache20,
	"ASL 2.0":                                            Apache20,
	"ASL, VERSION 2":                                     Apache20,
	"BERKELEY SOFTWARE DISTRIBUTION (BSD)":               BSD2Clause,
	"BOOST LICENSE V1.0":                                 BSL10,
	"BOOST SOFTWARE LICENSE 1.0 (BSL-1.0)":               BSL10,
	"BOOST SOFTWARE":                                     BSL10,
	"BOUNCY CASTLE":                                      MIT,
	"BSD (3-CLAUSE)":                                     BSD3Clause,
	"BSD - SEE NDG/HTTPSCLIENT/LICENSE FILE FOR DETAILS": BSD3Clause,
	"BSD 2 CLAUSE":                                       BSD2Clause,
	"BSD 2":                                              BSD2Clause,
	"BSD 2-CLAUSE":                                       BSD2Clause,
	"BSD 3 CLAUSE":                                       BSD3Clause,
	"BSD 3":                                              BSD3Clause,
	"BSD 3-CLAUSE NEW":                                   BSD3Clause,
	"BSD 3-CLAUSE \"NEW\" OR \"REVISED\" LICENSE (BSD-3-CLAUSE)": BSD3Clause,
	"BSD 3-CLAUSE":            BSD3Clause,
	"BSD 4 CLAUSE":            BSD4Clause,
	"BSD FOUR CLAUSE":         BSD4Clause,
	"BSD LICENSE 3":           BSD3Clause,
	"BSD LICENSE FOR HSQL":    BSD3Clause,
	"BSD NEW":                 BSD3Clause,
	"BSD THREE CLAUSE":        BSD3Clause,
	"BSD TWO CLAUSE":          BSD2Clause,
	"BSD-3 CLAUSE":            BSD3Clause,
	"BSD-STYLE + ATTRIBUTION": BSD3ClauseAttribution,
	"CC BY-NC-SA 2.0":         CCBYNCSA20,
	"CC BY-NC-SA 2.5":         CCBYNCSA25,
	"CC BY-NC-SA 3.0":         CCBYNCSA30,
	"CC BY-NC-SA 4.0":         CCBYNCSA40,
	"CC BY-SA 2.0":            CCBYSA20,
	"CC BY-SA 2.5":            CCBYSA25,
	"CC BY-SA 3.0":            CCBYSA30,
	"CC BY-SA 4.0":            CCBYSA40,
	"CC0 1.0 UNIVERSAL (CC0 1.0) PUBLIC DOMAIN DEDICATION": CC010,
	"CC0 1.0 UNIVERSAL": CC010,
	"CDDL 1.0":          CDDL10,
	"CDDL 1.1":          CDDL11,
	"CDDL V1.0":         CDDL10,
	"CDDL V1.1":         CDDL11,
	"CDDL, V1.0":        CDDL10,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE (CDDL) V1.0":                              CDDL10,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE (CDDL) V1.1":                              CDDL11,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE (CDDL) VERSION 1.0":                       CDDL10,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE (CDDL) VERSION 1.1":                       CDDL11,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE (CDDL), VERSION 1.1":                      CDDL11,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE 1.0 (CDDL-1.0)":                           CDDL10,
	"COMMON DEVELOPMENT AND DISTRIBUTION LICENSE 1.1 (CDDL-1.1)":                           CDDL11,
	"COMMON PUBLIC LICENSE - V 1.0":                                                        CPL10,
	"COMMON PUBLIC LICENSE VERSION 1.0":                                                    CPL10,
	"COMMON PUBLIC":                                                                        CPL10,
	"CPAL 1.0":                                                                             CPAL10,
	"CPAL V1.0":                                                                            CPAL10,
	"CREATIVE COMMONS - ATTRIBUTION 4.0 INTERNATIONAL":                                     CCBY40,
	"CREATIVE COMMONS 3.0 BY-SA":                                                           CCBYSA30,
	"CREATIVE COMMONS 3.0":                                                                 CCBY30,
	"CREATIVE COMMONS ATTRIBUTION 1.0":                                                     CCBY10,
	"CREATIVE COMMONS ATTRIBUTION 2.5":                                                     CCBY25,
	"CREATIVE COMMONS ATTRIBUTION 3.0 UNPORTED (CC BY 3.0)":                                CCBY30,
	"CREATIVE COMMONS ATTRIBUTION 3.0":                                                     CCBY30,
	"CREATIVE COMMONS ATTRIBUTION 4.0 INTERNATIONAL (CC BY 4.0)":                           CCBY40,
	"CREATIVE COMMONS ATTRIBUTION 4.0 INTERNATIONAL PUBLIC":                                CCBY40,
	"CREATIVE COMMONS ATTRIBUTION 4.0":                                                     CCBY40,
	"CREATIVE COMMONS ATTRIBUTION-NONCOMMERCIAL 4.0 INTERNATIONAL":                         CCBYNC40,
	"CREATIVE COMMONS ATTRIBUTION-NONCOMMERCIAL-NODERIVATIVES 4.0 INTERNATIONAL":           CCBYNCND40,
	"CREATIVE COMMONS ATTRIBUTION-NONCOMMERCIAL-SHAREALIKE 3.0 UNPORTED (CC BY-NC-SA 3.0)": CCBYNCSA30,
	"CREATIVE COMMONS ATTRIBUTION-NONCOMMERCIAL-SHAREALIKE 4.0 INTERNATIONAL PUBLIC":       CCBYNCSA40,
	"CREATIVE COMMONS CC0":                                                                 CC010,
	"CREATIVE COMMONS GNU LGPL, VERSION 2.1":                                               LGPL21,
	"CREATIVE COMMONS LICENSE ATTRIBUTION-NODERIVS 3.0 UNPORTED":                           CCBYNCND30,
	"CREATIVE COMMONS LICENSE ATTRIBUTION-NONCOMMERCIAL-SHAREALIKE 3.0 UNPORTED":           CCBYNCSA30,
	"CREATIVE COMMONS ZERO":                                                                CC010,
	"ECLIPSE 1.0":                                                                          EPL10,
	"ECLIPSE 2.0":                                                                          EPL20,
	"ECLIPSE DISTRIBUTION LICENSE (EDL), VERSION 1.0":                                      BSD3Clause,
	"ECLIPSE DISTRIBUTION LICENSE (NEW BSD LICENSE)":                                       BSD3Clause,
	"ECLIPSE DISTRIBUTION LICENSE - V 1.0":                                                 BSD3Clause,
	"ECLIPSE DISTRIBUTION LICENSE - VERSION 1.0":                                           BSD3Clause,
	"ECLIPSE DISTRIBUTION LICENSE V. 1.0":                                                  BSD3Clause,
	"ECLIPSE DISTRIBUTION LICENSE V1.0":                                                    BSD3Clause,
	"ECLIPSE PUBLIC LICENSE (EPL) 1.0":                                                     EPL10,
	"ECLIPSE PUBLIC LICENSE (EPL) 2.0":                                                     EPL20,
	"ECLIPSE PUBLIC LICENSE (EPL), VERSION 1.0":                                            EPL10,
	"ECLIPSE PUBLIC LICENSE - V 1.0":                                                       EPL10,
	"ECLIPSE PUBLIC LICENSE - V 2.0":                                                       EPL20,
	"ECLIPSE PUBLIC LICENSE - V1.0":                                                        EPL10,
	"ECLIPSE PUBLIC LICENSE - VERSION 1.0":                                                 EPL10,
	"ECLIPSE PUBLIC LICENSE - VERSION 2.0":                                                 EPL20,
	"ECLIPSE PUBLIC LICENSE 1.0 (EPL-1.0)":                                                 EPL10,
	"ECLIPSE PUBLIC LICENSE 1.0":                                                           EPL10,
	"ECLIPSE PUBLIC LICENSE 2.0 (EPL-2.0)":                                                 EPL20,
	"ECLIPSE PUBLIC LICENSE V. 2.0":                                                        EPL20,
	"ECLIPSE PUBLIC LICENSE V1.0":                                                          EPL10,
	"ECLIPSE PUBLIC LICENSE V2.0":                                                          EPL20,
	"ECLIPSE PUBLIC LICENSE VERSION 1.0":                                                   EPL10,
	"ECLIPSE PUBLIC LICENSE VERSION 2.0":                                                   EPL20,
	"ECLIPSE PUBLIC LICENSE, VERSION 1.0":                                                  EPL10,
	"ECLIPSE PUBLIC LICENSE, VERSION 2.0":                                                  EPL20,
	"ECLIPSE PUBLIC":                                                                       EPL10,
	"ECLIPSE PUBLISH LICENSE, VERSION 1.0":                                                 EPL10,
	"EDL 1.0":                                                                              BSD3Clause,
	"EPL (ECLIPSE PUBLIC LICENSE), V1.0 OR LATER":                                          EPL10,
	"EPL 1.0":                          EPL10,
	"EPL 2.0":                          EPL20,
	"EPL V1.0":                         EPL10,
	"EPL V2.0":                         EPL20,
	"EU PUBLIC LICENSE 1.0 (EUPL 1.0)": EUPL10,
	"EU PUBLIC LICENSE 1.1 (EUPL 1.1)": EUPL11,
	"EUPL 1.0":                         EUPL10,
	"EUPL 1.1":                         EUPL11,
	"EUPL V1.0":                        EUPL10,
	"EUPL V1.1":                        EUPL11,
	"EUROPEAN UNION PUBLIC LICENSE (EUPL V.1.1)":                                 EUPL11,
	"EUROPEAN UNION PUBLIC LICENSE 1.0 (EUPL 1.0)":                               EUPL10,
	"EUROPEAN UNION PUBLIC LICENSE 1.0":                                          EUPL10,
	"EUROPEAN UNION PUBLIC LICENSE 1.1 (EUPL 1.1)":                               EUPL11,
	"EUROPEAN UNION PUBLIC LICENSE 1.1":                                          EUPL11,
	"EUROPEAN UNION PUBLIC LICENSE, VERSION 1.1":                                 EUPL11,
	"EXPAT (MIT/X11)":                                                            MIT,
	"GENERAL PUBLIC LICENSE 2.0 (GPL)":                                           GPL20,
	"GNU AFFERO GENERAL PUBLIC LICENSE V3 (AGPL-3.0)":                            AGPL30,
	"GNU AFFERO GENERAL PUBLIC LICENSE V3 (AGPLV3)":                              AGPL30,
	"GNU AFFERO GENERAL PUBLIC LICENSE V3 OR LATER (AGPL3+)":                     AGPL30,
	"GNU AFFERO GENERAL PUBLIC LICENSE V3 OR LATER (AGPLV3+)":                    AGPL30,
	"GNU AFFERO GENERAL PUBLIC LICENSE V3":                                       AGPL30,
	"GNU AFFERO GENERAL PUBLIC LICENSE, VERSION 3":                               AGPL30,
	"GNU FREE DOCUMENTATION LICENSE (GFDL-1.3)":                                  GFDL13,
	"GNU GENERAL LESSER PUBLIC LICENSE (LGPL) VERSION 2.1":                       LGPL21,
	"GNU GENERAL LESSER PUBLIC LICENSE (LGPL) VERSION 3.0":                       LGPL30,
	"GNU GENERAL PUBLIC LIBRARY":                                                 GPL30,
	"GNU GENERAL PUBLIC LICENSE (GPL) V. 2":                                      GPL20,
	"GNU GENERAL PUBLIC LICENSE (GPL) V. 3":                                      GPL30,
	"GNU GENERAL PUBLIC LICENSE (GPL), VERSION 2, WITH CLASSPATH EXCEPTION":      GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE (GPL), VERSION 2, WITH THE CLASSPATH EXCEPTION":  GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE 3":                                               GPL30,
	"GNU GENERAL PUBLIC LICENSE V2 (GPLV2)":                                      GPL20,
	"GNU GENERAL PUBLIC LICENSE V2 OR LATER (GPLV2+)":                            GPL20,
	"GNU GENERAL PUBLIC LICENSE V2.0 ONLY, WITH CLASSPATH EXCEPTION":             GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE V3 (GPLV3)":                                      GPL30,
	"GNU GENERAL PUBLIC LICENSE V3 OR LATER (GPLV3+)":                            GPL30,
	"GNU GENERAL PUBLIC LICENSE VERSION 2 (GPLV2)":                               GPL20,
	"GNU GENERAL PUBLIC LICENSE VERSION 2":                                       GPL20,
	"GNU GENERAL PUBLIC LICENSE VERSION 2, JUNE 1991":                            GPL20,
	"GNU GENERAL PUBLIC LICENSE VERSION 3 (GPL V3)":                              GPL30,
	"GNU GENERAL PUBLIC LICENSE, VERSION 2 (GPL2), WITH THE CLASSPATH EXCEPTION": GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE, VERSION 2 WITH THE CLASSPATH EXCEPTION":         GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE, VERSION 2 WITH THE GNU CLASSPATH EXCEPTION":     GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE, VERSION 2":                                      GPL20,
	"GNU GENERAL PUBLIC LICENSE, VERSION 2, WITH THE CLASSPATH EXCEPTION":        GPL20withclasspathexception,
	"GNU GENERAL PUBLIC LICENSE, VERSION 3":                                      GPL30,
	"GNU GPL V2":                                                                 GPL20,
	"GNU GPL V3":                                                                 GPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE (LGPL), VERSION 2.1 OR LATER":             LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE (LGPL), VERSION 2.1":                      LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE (LGPL), VERSION 3":                        LGPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE 2.1":                                      LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE V2 (LGPLV2)":                              LGPL20,
	"GNU LESSER GENERAL PUBLIC LICENSE V2 OR LATER (LGPLV2+)":                    LGPL20,
	"GNU LESSER GENERAL PUBLIC LICENSE V3 (LGPLV3)":                              LGPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE V3 OR LATER (LGPLV3+)":                    LGPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE V3":                                       LGPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE V3.0":                                     LGPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE VERSION 2.1 (LGPLV2.1)":                   LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE VERSION 2.1":                              LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE VERSION 2.1, FEBRUARY 1999":               LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE, VERSION 2.1":                             LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE, VERSION 2.1, FEBRUARY 1999":              LGPL21,
	"GNU LESSER GENERAL PUBLIC LICENSE, VERSION 3":                               LGPL30,
	"GNU LESSER GENERAL PUBLIC LICENSE, VERSION 3.0":                             LGPL30,
	"GNU LGP (GNU GENERAL PUBLIC LICENSE), V2 OR LATER":                          LGPL20,
	"GNU LGPL (GNU LESSER GENERAL PUBLIC LICENSE), V2.1 OR LATER":                LGPL21,
	"GNU LGPL 2":    LGPL20,
	"GNU LGPL 2.1":  LGPL21,
	"GNU LGPL 3.0":  LGPL30,
	"GNU LGPL V2":   LGPL20,
	"GNU LGPL V2.1": LGPL21,
	"GNU LGPL V3":   LGPL30,
	"GNU LIBRARY GENERAL PUBLIC LICENSE V2.1 OR LATER":                  LGPL21,
	"GNU LIBRARY OR LESSER GENERAL PUBLIC LICENSE VERSION 2.0 (LGPLV2)": LGPL20,
	"GPL (≥ 3)":                       GPL30,
	"GPL 2 WITH CLASSPATH EXCEPTION":  GPL20withclasspathexception,
	"GPL 2":                           GPL20,
	"GPL 3":                           GPL30,
	"GPL V2 WITH CLASSPATH EXCEPTION": GPL20withclasspathexception,
	"GPL V2":                          GPL20,
	"GPL V3":                          GPL30,
	"GPL VERSION 2":                   GPL20,
	"GPL-2+ WITH AUTOCONF EXCEPTION":  GPL20withautoconfexception,
	"GPL-3+ WITH AUTOCONF EXCEPTION":  GPL30withautoconfexception,
	"GPL2 W/ CPE":                     GPL20withclasspathexception,
	"GPLV2 LICENSE, INCLUDES THE CLASSPATH EXCEPTION":                       GPL20withclasspathexception,
	"GPLV2 WITH CLASSPATH EXCEPTION":                                        GPL20withclasspathexception,
	"HSQLDB LICENSE, A BSD OPEN SOURCE":                                     BSD3Clause,
	"HTTP://ANT-CONTRIB.SOURCEFORGE.NET/TASKS/LICENSE.TXT":                  Apache11,
	"HTTP://ASM.OW2.ORG/LICENSE.HTML":                                       BSD3Clause,
	"HTTP://CREATIVECOMMONS.ORG/PUBLICDOMAIN/ZERO/1.0/LEGALCODE":            CC010,
	"HTTP://EN.WIKIPEDIA.ORG/WIKI/ZLIB_LICENSE":                             Zlib,
	"HTTP://JSON.CODEPLEX.COM/LICENSE":                                      MIT,
	"HTTP://POLYMER.GITHUB.IO/LICENSE.TXT":                                  BSD3Clause,
	"HTTP://WWW.APACHE.ORG/LICENSES/LICENSE-2.0":                            Apache20,
	"HTTP://WWW.APACHE.ORG/LICENSES/LICENSE-2.0.HTML":                       Apache20,
	"HTTP://WWW.APACHE.ORG/LICENSES/LICENSE-2.0.TXT":                        Apache20,
	"HTTP://WWW.GNU.ORG/COPYLEFT/LESSER.HTML":                               LGPL30,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-ND/1.0":                     CCBYNCND10,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-ND/2.0":                     CCBYNCND20,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-ND/2.5":                     CCBYNCND25,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-ND/3.0":                     CCBYNCND30,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-ND/4.0":                     CCBYNCND40,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-SA/1.0":                     CCBYNCSA10,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-SA/2.0":                     CCBYNCSA20,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-SA/2.5":                     CCBYNCSA25,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-SA/3.0":                     CCBYNCSA30,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-NC-SA/4.0":                     CCBYNCSA40,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-ND/1.0":                        CCBYND10,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-ND/2.0":                        CCBYND20,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-ND/2.5":                        CCBYND25,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-ND/3.0":                        CCBYND30,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-ND/4.0":                        CCBYND40,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-SA/1.0":                        CCBYSA10,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-SA/2.0":                        CCBYSA20,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-SA/2.5":                        CCBYSA25,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-SA/3.0":                        CCBYSA30,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY-SA/4.0":                        CCBYSA40,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY/1.0":                           CCBY10,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY/2.0":                           CCBY20,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY/2.5":                           CCBY25,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY/3.0":                           CCBY30,
	"HTTPS://CREATIVECOMMONS.ORG/LICENSES/BY/4.0":                           CCBY40,
	"HTTPS://CREATIVECOMMONS.ORG/PUBLICDOMAIN/ZERO/1.0/":                    CC010,
	"HTTPS://GITHUB.COM/DOTNET/CORE-SETUP/BLOB/MASTER/LICENSE.TXT":          MIT,
	"HTTPS://GITHUB.COM/DOTNET/COREFX/BLOB/MASTER/LICENSE.TXT":              MIT,
	"HTTPS://RAW.GITHUB.COM/RDFLIB/RDFLIB/MASTER/LICENSE":                   BSD3Clause,
	"HTTPS://RAW.GITHUBUSERCONTENT.COM/ASPNET/ASPNETCORE/2.0.0/LICENSE.TXT": Apache20,
	"HTTPS://RAW.GITHUBUSERCONTENT.COM/ASPNET/HOME/2.0.0/LICENSE.TXT":       Apache20,
	"HTTPS://RAW.GITHUBUSERCONTENT.COM/NUGET/NUGET.CLIENT/DEV/LICENSE.TXT":  Apache20,
	"HTTPS://WWW.APACHE.ORG/LICENSES/LICENSE-2.0":                           Apache20,
	"HTTPS://WWW.ECLIPSE.ORG/LEGAL/EPL-V10.HTML":                            EPL10,
	"HTTPS://WWW.ECLIPSE.ORG/LEGAL/EPL-V20.HTML":                            EPL20,
	"IBM PUBLIC":         IPL10,
	"ISC LICENSE (ISCL)": ISC,
	"JYTHON SOFTWARE":    Python20,
	"KIRKK.COM BSD":      BSD3Clause,
	"LESSER GENERAL PUBLIC LICENSE, VERSION 3 OR GREATER": LGPL30,
	"LGPL 2":              LGPL20,
	"LGPL 2.0":            LGPL20,
	"LGPL 2.1":            LGPL21,
	"LGPL 3":              LGPL30,
	"LGPL 3.0":            LGPL30,
	"LGPL V2":             LGPL20,
	"LGPL V2.1":           LGPL21,
	"LGPL V3":             LGPL30,
	"LGPL V3.0":           LGPL30,
	"LGPL, V2.1 OR LATER": LGPL21,
	"LGPL, VERSION 2.1":   LGPL21,
	"LGPL, VERSION 3.0":   LGPL30,
	"LGPLV3 OR LATER":     LGPL30,
	"LICENSE AGREEMENT FOR OPEN SOURCE COMPUTER VISION LIBRARY (3-CLAUSE BSD LICENSE)": BSD3Clause,
	"MIT (HTTP://MOOTOOLS.NET/LICENSE.TXT)":                                            MIT,
	"MIT / HTTP://REM.MIT-LICENSE.ORG":                                                 MIT,
	"MIT LICENSE (HTTP://OPENSOURCE.ORG/LICENSES/MIT)":                                 MIT,
	"MIT LICENSE (MIT)": MIT,
	"MIT LICENSE(MIT)":  MIT,
	"MIT LICENSED. HTTP://WWW.OPENSOURCE.ORG/LICENSES/MIT-LICENSE.PHP": MIT,
	"MIT/EXPAT": MIT,
	"MOCKRUNNER LICENSE, BASED ON APACHE SOFTWARE LICENSE, VERSION 1.1": Apache11,
	"MODIFIED BSD":                                  BSD3Clause,
	"MOZILLA PUBLIC LICENSE 1.0 (MPL)":              MPL10,
	"MOZILLA PUBLIC LICENSE 1.1 (MPL 1.1)":          MPL11,
	"MOZILLA PUBLIC LICENSE 2.0 (MPL 2.0)":          MPL20,
	"MOZILLA PUBLIC LICENSE V 2.0":                  MPL20,
	"MOZILLA PUBLIC LICENSE VERSION 1.0":            MPL10,
	"MOZILLA PUBLIC LICENSE VERSION 1.1":            MPL11,
	"MOZILLA PUBLIC LICENSE VERSION 2.0":            MPL20,
	"MOZILLA PUBLIC LICENSE, VERSION 2.0":           MPL20,
	"MPL 1":                                         MPL10,
	"MPL 1.0":                                       MPL10,
	"MPL 1.1":                                       MPL11,
	"MPL 2":                                         MPL20,
	"MPL 2.0":                                       MPL20,
	"MPL V2":                                        MPL20,
	"NCSA OPEN SOURCE":                              NCSA,
	"NETSCAPE PUBLIC LICENSE (NPL)":                 NPL10,
	"NETSCAPE PUBLIC":                               NPL10,
	"NEW BSD":                                       BSD3Clause,
	"OPEN SOFTWARE LICENSE 3.0 (OSL-3.0)":           OSL30,
	"OPEN SOFTWARE LICENSE V. 3.0":                  OSL30,
	"PERL ARTISTIC V2":                              Artistic10Perl,
	"PUBLIC DOMAIN":                                 Unlicense,
	"PUBLIC DOMAIN (CC0-1.0)":                       CC010,
	"PUBLIC DOMAIN, PER CREATIVE COMMONS CC0":       CC010,
	"QT PUBLIC LICENSE (QPL)":                       QPL10,
	"QT PUBLIC":                                     QPL10,
	"REVISED BSD":                                   BSD3Clause,
	"RUBY'S":                                        Ruby,
	"SEQUENCE LIBRARY LICENSE (BSD-LIKE)":           BSD3Clause,
	"SIL OPEN FONT LICENSE 1.1 (OFL-1.1)":           OFL11,
	"SIL OPEN FONT LICENSE VERSION 1.1":             OFL11,
	"SIMPLIFIED BSD LISCENCE":                       BSD2Clause,
	"SIMPLIFIED BSD":                                BSD2Clause,
	"SUN INDUSTRY STANDARDS SOURCE LICENSE (SISSL)": SISSL,
	"THREE-CLAUSE BSD-STYLE":                        BSD3Clause,
	"TWO-CLAUSE BSD-STYLE":                          BSD2Clause,
	"UNIVERSAL PERMISSIVE LICENSE (UPL)":            UPL10,
	"UNIVERSAL PERMISSIVE LICENSE, VERSION 1.0":     UPL10,
	"UNLICENSE (UNLICENSE)":                         Unlicense,
	"W3C SOFTWARE":                                  W3C,
	"ZLIB / LIBPNG":                                 ZlibAcknowledgement,
	"ZLIB/LIBPNG":                                   ZlibAcknowledgement,
	"ZOPE 1.1":                                      ZPL11,
	"ZOPE 2.0":                                      ZPL20,
	"ZOPE 2.1":                                      ZPL21,
	"ZOPE PUBLIC":                                   ZPL21,
	"ZOPE V2.1":                                     ZPL21,
	"ZPL 2.1":                                       ZPL21,
	"['MIT']":                                       MIT,
}

// pythonLicenseExceptions contains licenses that we cannot separate correctly using our logic.
// first word after separator (or/and) => license name
var pythonLicenseExceptions = map[string]string{
	"lesser":       "GNU Library or Lesser General Public License (LGPL)",
	"distribution": "Common Development and Distribution License 1.0 (CDDL-1.0)",
	"disclaimer":   "Historical Permission Notice and Disclaimer (HPND)",
}

// Split licenses without considering "and"/"or"
// examples:
// 'GPL-1+,GPL-2' => {"GPL-1+", "GPL-2"}
// 'GPL-1+ or Artistic or Artistic-dist' => {"GPL-1+", "Artistic", "Artistic-dist"}
// 'LGPLv3+_or_GPLv2+' => {"LGPLv3+", "GPLv2"}
// 'BSD-3-CLAUSE and GPL-2' => {"BSD-3-CLAUSE", "GPL-2"}
// 'GPL-1+ or Artistic, and BSD-4-clause-POWERDOG' => {"GPL-1+", "Artistic", "BSD-4-clause-POWERDOG"}
// 'BSD 3-Clause License or Apache License, Version 2.0' => {"BSD 3-Clause License", "Apache License, Version 2.0"}
// var LicenseSplitRegexp = regexp.MustCompile("(,?[_ ]+or[_ ]+)|(,?[_ ]+and[_ ])|(,[ ]*)")

var licenseSplitRegexp = regexp.MustCompile("(,?[_ ]+(?:or|and)[_ ]+)|(,[ ]*)")

func Normalize(name string) string {
	// standardize space, including newline
	license := strings.Join(strings.Fields(name), " ")
	license = strings.TrimSpace(license)
	license = strings.ToUpper(license)
	license = strings.ReplaceAll(license, "LICENCE", "LICENSE")
	license = strings.TrimPrefix(license, "THE ")
	license = strings.TrimSuffix(license, " LICENSE")
	license = strings.TrimSuffix(license, " LICENSED")
	license = strings.TrimSuffix(license, "-LICENSE")
	license = strings.TrimSuffix(license, "-LICENSED")
	// suffixes from https://spdx.dev/learn/handling-license-info/
	license = strings.TrimSuffix(license, "+")
	// Note: -only and -or-later GNU licenses could also be matched with new SPDX ids such as GPL-3.0-only
	// if they are added to category.go, but those new ids are not supported by google/licenseclassifier
	license = strings.TrimSuffix(license, "-ONLY")
	license = strings.TrimSuffix(license, "-OR-LATER")
	if l, ok := mapping[license]; ok {
		return l
	}
	return name
}

func SplitLicenses(str string) []string {
	if str == "" {
		return nil
	}
	var licenses []string
	for _, maybeLic := range licenseSplitRegexp.Split(str, -1) {
		lower := strings.ToLower(maybeLic)
		firstWord, _, _ := strings.Cut(lower, " ")
		if len(licenses) > 0 {
			// e.g. `Apache License, Version 2.0`
			if firstWord == "ver" || firstWord == "version" {
				licenses[len(licenses)-1] += ", " + maybeLic
				continue
				// e.g. `GNU Lesser General Public License v2 or later (LGPLv2+)`
			} else if firstWord == "later" {
				licenses[len(licenses)-1] += " or " + maybeLic
				continue
			} else if lic, ok := pythonLicenseExceptions[firstWord]; ok {
				// Check `or` and `and` separators
				if lic == licenses[len(licenses)-1]+" or "+maybeLic || lic == licenses[len(licenses)-1]+" and "+maybeLic {
					licenses[len(licenses)-1] = lic
				}
				continue
			}
		}
		licenses = append(licenses, maybeLic)
	}
	return licenses
}
