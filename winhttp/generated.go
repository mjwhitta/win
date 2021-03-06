// Code generated by tools/defines.go; DO NOT EDIT.
package winhttp

const (
	InternetDefaultPort                             uintptr = 0
	InternetDefaultHTTPPort                         uintptr = 80
	InternetDefaultHTTPsPort                        uintptr = 443
	InternetSchemeHTTP                              uintptr = 1
	InternetSchemeHTTPs                             uintptr = 2
	InternetSchemeFtp                               uintptr = 3
	InternetSchemeSocks                             uintptr = 4
	IcuEscape                                       uintptr = 0x80000000
	IcuEscapeAuthority                              uintptr = 0x00002000
	IcuRejectUserpwd                                uintptr = 0x00004000
	WinhttpFlagAsync                                uintptr = 0x10000000
	WinhttpFlagEscapePercent                        uintptr = 0x00000004
	WinhttpFlagNullCodepage                         uintptr = 0x00000008
	WinhttpFlagEscapeDisable                        uintptr = 0x00000040
	WinhttpFlagEscapeDisableQuery                   uintptr = 0x00000080
	WinhttpFlagBypassProxyCache                     uintptr = 0x00000100
	WinhttpFlagRefresh                              uintptr = WinhttpFlagBypassProxyCache
	WinhttpFlagSecure                               uintptr = 0x00800000
	WinhttpAccessTypeDefaultProxy                   uintptr = 0
	WinhttpAccessTypeNoProxy                        uintptr = 1
	WinhttpAccessTypeNamedProxy                     uintptr = 3
	WinhttpAccessTypeAutomaticProxy                 uintptr = 4
	WinhttpNoProxyName                              uintptr = 0
	WinhttpNoProxyBypass                            uintptr = 0
	WinhttpNoClientCertContext                      uintptr = 0
	WinhttpNoReferer                                uintptr = 0
	WinhttpDefaultAcceptTypes                       uintptr = 0
	WinhttpNoAdditionalHeaders                      uintptr = 0
	WinhttpNoRequestData                            uintptr = 0
	WinhttpHeaderNameByIndex                        uintptr = 0
	WinhttpNoOutputBuffer                           uintptr = 0
	WinhttpNoHeaderIndex                            uintptr = 0
	WinhttpAddreqIndexMask                          uintptr = 0x0000FFFF
	WinhttpAddreqFlagsMask                          uintptr = 0xFFFF0000
	WinhttpAddreqFlagAddIfNew                       uintptr = 0x10000000
	WinhttpAddreqFlagAdd                            uintptr = 0x20000000
	WinhttpAddreqFlagCoalesceWithComma              uintptr = 0x40000000
	WinhttpAddreqFlagCoalesceWithSemicolon          uintptr = 0x01000000
	WinhttpAddreqFlagCoalesce                       uintptr = WinhttpAddreqFlagCoalesceWithComma
	WinhttpAddreqFlagReplace                        uintptr = 0x80000000
	WinhttpIgnoreRequestTotalLength                 uintptr = 0
	WinhttpFirstOption                              uintptr = WinhttpOptionCallback
	WinhttpOptionCallback                           uintptr = 1
	WinhttpOptionResolveTimeout                     uintptr = 2
	WinhttpOptionConnectTimeout                     uintptr = 3
	WinhttpOptionConnectRetries                     uintptr = 4
	WinhttpOptionSendTimeout                        uintptr = 5
	WinhttpOptionReceiveTimeout                     uintptr = 6
	WinhttpOptionReceiveResponseTimeout             uintptr = 7
	WinhttpOptionHandleType                         uintptr = 9
	WinhttpOptionReadBufferSize                     uintptr = 12
	WinhttpOptionWriteBufferSize                    uintptr = 13
	WinhttpOptionParentHandle                       uintptr = 21
	WinhttpOptionExtendedError                      uintptr = 24
	WinhttpOptionSecurityFlags                      uintptr = 31
	WinhttpOptionSecurityCertificateStruct          uintptr = 32
	WinhttpOptionUrl                                uintptr = 34
	WinhttpOptionSecurityKeyBitness                 uintptr = 36
	WinhttpOptionProxy                              uintptr = 38
	WinhttpOptionProxyResultEntry                   uintptr = 39
	WinhttpOptionUserAgent                          uintptr = 41
	WinhttpOptionContextValue                       uintptr = 45
	WinhttpOptionClientCertContext                  uintptr = 47
	WinhttpOptionRequestPriority                    uintptr = 58
	WinhttpOptionHTTPVersion                        uintptr = 59
	WinhttpOptionDisableFeature                     uintptr = 63
	WinhttpOptionCodepage                           uintptr = 68
	WinhttpOptionMaxConnsPerServer                  uintptr = 73
	WinhttpOptionMaxConnsPer10Server                uintptr = 74
	WinhttpOptionAutologonPolicy                    uintptr = 77
	WinhttpOptionServerCertContext                  uintptr = 78
	WinhttpOptionEnableFeature                      uintptr = 79
	WinhttpOptionWorkerThreadCount                  uintptr = 80
	WinhttpOptionPassportCobrandingText             uintptr = 81
	WinhttpOptionPassportCobrandingUrl              uintptr = 82
	WinhttpOptionConfigurePassportAuth              uintptr = 83
	WinhttpOptionSecureProtocols                    uintptr = 84
	WinhttpOptionEnabletracing                      uintptr = 85
	WinhttpOptionPassportSignOut                    uintptr = 86
	WinhttpOptionPassportReturnUrl                  uintptr = 87
	WinhttpOptionRedirectPolicy                     uintptr = 88
	WinhttpOptionMaxHTTPAutomaticRedirects          uintptr = 89
	WinhttpOptionMaxHTTPStatusContinue              uintptr = 90
	WinhttpOptionMaxResponseHeaderSize              uintptr = 91
	WinhttpOptionMaxResponseDrainSize               uintptr = 92
	WinhttpOptionConnectionInfo                     uintptr = 93
	WinhttpOptionClientCertIssuerList               uintptr = 94
	WinhttpOptionSpn                                uintptr = 96
	WinhttpOptionGlobalProxyCreds                   uintptr = 97
	WinhttpOptionGlobalServerCreds                  uintptr = 98
	WinhttpOptionUnloadNotifyEvent                  uintptr = 99
	WinhttpOptionRejectUserpwdInUrl                 uintptr = 100
	WinhttpOptionUseGlobalServerCredentials         uintptr = 101
	WinhttpOptionReceiveProxyConnectResponse        uintptr = 103
	WinhttpOptionIsProxyConnectResponse             uintptr = 104
	WinhttpOptionServerSpnUsed                      uintptr = 106
	WinhttpOptionProxySpnUsed                       uintptr = 107
	WinhttpOptionServerCbt                          uintptr = 108
	WinhttpOptionUnsafeHeaderParsing                uintptr = 110
	WinhttpOptionAssuredNonBlockingCallbacks        uintptr = 111
	WinhttpOptionUpgradeToWebSocket                 uintptr = 114
	WinhttpOptionWebSocketCloseTimeout              uintptr = 115
	WinhttpOptionWebSocketKeepaliveInterval         uintptr = 116
	WinhttpOptionDecompression                      uintptr = 118
	WinhttpOptionWebSocketReceiveBufferSize         uintptr = 122
	WinhttpOptionWebSocketSendBufferSize            uintptr = 123
	WinhttpOptionTcpPriorityHint                    uintptr = 128
	WinhttpOptionConnectionFilter                   uintptr = 131
	WinhttpOptionEnableHTTPProtocol                 uintptr = 133
	WinhttpOptionHTTPProtocolUsed                   uintptr = 134
	WinhttpOptionKdcProxySettings                   uintptr = 136
	WinhttpOptionEncodeExtra                        uintptr = 138
	WinhttpOptionDisableStreamQueue                 uintptr = 139
	WinhttpOptionIpv6FastFallback                   uintptr = 140
	WinhttpOptionConnectionStatsV0                  uintptr = 141
	WinhttpOptionRequestTimes                       uintptr = 142
	WinhttpOptionExpireConnection                   uintptr = 143
	WinhttpOptionDisableSecureProtocolFallback      uintptr = 144
	WinhttpOptionHTTPProtocolRequired               uintptr = 145
	WinhttpOptionRequestStats                       uintptr = 146
	WinhttpOptionServerCertChainContext             uintptr = 147
	WinhttpLastOption                               uintptr = WinhttpOptionServerCertChainContext
	WinhttpOptionUsername                           uintptr = 0x1000
	WinhttpOptionPassword                           uintptr = 0x1001
	WinhttpOptionProxyUsername                      uintptr = 0x1002
	WinhttpOptionProxyPassword                      uintptr = 0x1003
	WinhttpConnsPerServerUnlimited                  uintptr = 0xFFFFFFFF
	WinhttpDecompressionFlagGzip                    uintptr = 0x00000001
	WinhttpDecompressionFlagDeflate                 uintptr = 0x00000002
	WinhttpDecompressionFlagAll                     uintptr = (WinhttpDecompressionFlagGzip | WinhttpDecompressionFlagDeflate)
	WinhttpProtocolFlagHTTP2                        uintptr = 0x1
	WinhttpProtocolMask                             uintptr = WinhttpProtocolFlagHTTP2
	WinhttpAutologonSecurityLevelMedium             uintptr = 0
	WinhttpAutologonSecurityLevelLow                uintptr = 1
	WinhttpAutologonSecurityLevelHigh               uintptr = 2
	WinhttpAutologonSecurityLevelDefault            uintptr = WinhttpAutologonSecurityLevelMedium
	WinhttpOptionRedirectPolicyNever                uintptr = 0
	WinhttpOptionRedirectPolicyDisallowHTTPsToHTTP  uintptr = 1
	WinhttpOptionRedirectPolicyAlways               uintptr = 2
	WinhttpOptionRedirectPolicyLast                 uintptr = WinhttpOptionRedirectPolicyAlways
	WinhttpOptionRedirectPolicyDefault              uintptr = WinhttpOptionRedirectPolicyDisallowHTTPsToHTTP
	WinhttpDisablePassportAuth                      uintptr = 0x00000000
	WinhttpEnablePassportAuth                       uintptr = 0x10000000
	WinhttpDisablePassportKeyring                   uintptr = 0x20000000
	WinhttpEnablePassportKeyring                    uintptr = 0x40000000
	WinhttpDisableCookies                           uintptr = 0x00000001
	WinhttpDisableRedirects                         uintptr = 0x00000002
	WinhttpDisableAuthentication                    uintptr = 0x00000004
	WinhttpDisableKeepAlive                         uintptr = 0x00000008
	WinhttpEnableSslRevocation                      uintptr = 0x00000001
	WinhttpEnableSslRevertImpersonation             uintptr = 0x00000002
	WinhttpDisableSpnServerPort                     uintptr = 0x00000000
	WinhttpEnableSpnServerPort                      uintptr = 0x00000001
	WinhttpOptionSpnMask                            uintptr = WinhttpEnableSpnServerPort
	WinhttpErrorBase                                uintptr = 12000
	ErrorWinhttpOutOfHandles                        uintptr = (WinhttpErrorBase + 1)
	ErrorWinhttpTimeout                             uintptr = (WinhttpErrorBase + 2)
	ErrorWinhttpInternalError                       uintptr = (WinhttpErrorBase + 4)
	ErrorWinhttpInvalidUrl                          uintptr = (WinhttpErrorBase + 5)
	ErrorWinhttpUnrecognizedScheme                  uintptr = (WinhttpErrorBase + 6)
	ErrorWinhttpNameNotResolved                     uintptr = (WinhttpErrorBase + 7)
	ErrorWinhttpInvalidOption                       uintptr = (WinhttpErrorBase + 9)
	ErrorWinhttpOptionNotSettable                   uintptr = (WinhttpErrorBase + 11)
	ErrorWinhttpShutdown                            uintptr = (WinhttpErrorBase + 12)
	ErrorWinhttpLoginFailure                        uintptr = (WinhttpErrorBase + 15)
	ErrorWinhttpOperationCancelled                  uintptr = (WinhttpErrorBase + 17)
	ErrorWinhttpIncorrectHandleType                 uintptr = (WinhttpErrorBase + 18)
	ErrorWinhttpIncorrectHandleState                uintptr = (WinhttpErrorBase + 19)
	ErrorWinhttpCannotConnect                       uintptr = (WinhttpErrorBase + 29)
	ErrorWinhttpConnectionError                     uintptr = (WinhttpErrorBase + 30)
	ErrorWinhttpResendRequest                       uintptr = (WinhttpErrorBase + 32)
	ErrorWinhttpSecureCertDateInvalid               uintptr = (WinhttpErrorBase + 37)
	ErrorWinhttpSecureCertCnInvalid                 uintptr = (WinhttpErrorBase + 38)
	ErrorWinhttpClientAuthCertNeeded                uintptr = (WinhttpErrorBase + 44)
	ErrorWinhttpSecureInvalidCa                     uintptr = (WinhttpErrorBase + 45)
	ErrorWinhttpSecureCertRevFailed                 uintptr = (WinhttpErrorBase + 57)
	ErrorWinhttpCannotCallBeforeOpen                uintptr = (WinhttpErrorBase + 100)
	ErrorWinhttpCannotCallBeforeSend                uintptr = (WinhttpErrorBase + 101)
	ErrorWinhttpCannotCallAfterSend                 uintptr = (WinhttpErrorBase + 102)
	ErrorWinhttpCannotCallAfterOpen                 uintptr = (WinhttpErrorBase + 103)
	ErrorWinhttpHeaderNotFound                      uintptr = (WinhttpErrorBase + 150)
	ErrorWinhttpInvalidServerResponse               uintptr = (WinhttpErrorBase + 152)
	ErrorWinhttpInvalidHeader                       uintptr = (WinhttpErrorBase + 153)
	ErrorWinhttpInvalidQueryRequest                 uintptr = (WinhttpErrorBase + 154)
	ErrorWinhttpHeaderAlreadyExists                 uintptr = (WinhttpErrorBase + 155)
	ErrorWinhttpRedirectFailed                      uintptr = (WinhttpErrorBase + 156)
	ErrorWinhttpSecureChannelError                  uintptr = (WinhttpErrorBase + 157)
	ErrorWinhttpBadAutoProxyScript                  uintptr = (WinhttpErrorBase + 166)
	ErrorWinhttpUnableToDownloadScript              uintptr = (WinhttpErrorBase + 167)
	ErrorWinhttpSecureInvalidCert                   uintptr = (WinhttpErrorBase + 169)
	ErrorWinhttpSecureCertRevoked                   uintptr = (WinhttpErrorBase + 170)
	ErrorWinhttpNotInitialized                      uintptr = (WinhttpErrorBase + 172)
	ErrorWinhttpSecureFailure                       uintptr = (WinhttpErrorBase + 175)
	ErrorWinhttpUnhandledScriptType                 uintptr = (WinhttpErrorBase + 176)
	ErrorWinhttpScriptExecutionError                uintptr = (WinhttpErrorBase + 177)
	ErrorWinhttpAutoProxyServiceError               uintptr = (WinhttpErrorBase + 178)
	ErrorWinhttpSecureCertWrongUsage                uintptr = (WinhttpErrorBase + 179)
	ErrorWinhttpAutodetectionFailed                 uintptr = (WinhttpErrorBase + 180)
	ErrorWinhttpHeaderCountExceeded                 uintptr = (WinhttpErrorBase + 181)
	ErrorWinhttpHeaderSizeOverflow                  uintptr = (WinhttpErrorBase + 182)
	ErrorWinhttpChunkedEncodingHeaderSizeOverflow   uintptr = (WinhttpErrorBase + 183)
	ErrorWinhttpResponseDrainOverflow               uintptr = (WinhttpErrorBase + 184)
	ErrorWinhttpClientCertNoPrivateKey              uintptr = (WinhttpErrorBase + 185)
	ErrorWinhttpClientCertNoAccessPrivateKey        uintptr = (WinhttpErrorBase + 186)
	ErrorWinhttpClientAuthCertNeededProxy           uintptr = (WinhttpErrorBase + 187)
	ErrorWinhttpSecureFailureProxy                  uintptr = (WinhttpErrorBase + 188)
	ErrorWinhttpReserved189                         uintptr = (WinhttpErrorBase + 189)
	ErrorWinhttpHTTPProtocolMismatch                uintptr = (WinhttpErrorBase + 190)
	WinhttpErrorLast                                uintptr = (WinhttpErrorBase + 188)
	WinhttpResetState                               uintptr = 0x00000001
	WinhttpResetSwpadCurrentNetwork                 uintptr = 0x00000002
	WinhttpResetSwpadAll                            uintptr = 0x00000004
	WinhttpResetScriptCache                         uintptr = 0x00000008
	WinhttpResetAll                                 uintptr = 0x0000FFFF
	WinhttpResetNotifyNetworkChanged                uintptr = 0x00010000
	WinhttpResetOutOfProc                           uintptr = 0x00020000
	WinhttpResetDiscardResolvers                    uintptr = 0x00040000
	HTTPStatusContinue                              uintptr = 100
	HTTPStatusSwitchProtocols                       uintptr = 101
	HTTPStatusOk                                    uintptr = 200
	HTTPStatusCreated                               uintptr = 201
	HTTPStatusAccepted                              uintptr = 202
	HTTPStatusPartial                               uintptr = 203
	HTTPStatusNoContent                             uintptr = 204
	HTTPStatusResetContent                          uintptr = 205
	HTTPStatusPartialContent                        uintptr = 206
	HTTPStatusWebdavMultiStatus                     uintptr = 207
	HTTPStatusAmbiguous                             uintptr = 300
	HTTPStatusMoved                                 uintptr = 301
	HTTPStatusRedirect                              uintptr = 302
	HTTPStatusRedirectMethod                        uintptr = 303
	HTTPStatusNotModified                           uintptr = 304
	HTTPStatusUseProxy                              uintptr = 305
	HTTPStatusRedirectKeepVerb                      uintptr = 307
	HTTPStatusPermanentRedirect                     uintptr = 308
	HTTPStatusBadRequest                            uintptr = 400
	HTTPStatusDenied                                uintptr = 401
	HTTPStatusPaymentReq                            uintptr = 402
	HTTPStatusForbidden                             uintptr = 403
	HTTPStatusNotFound                              uintptr = 404
	HTTPStatusBadMethod                             uintptr = 405
	HTTPStatusNoneAcceptable                        uintptr = 406
	HTTPStatusProxyAuthReq                          uintptr = 407
	HTTPStatusRequestTimeout                        uintptr = 408
	HTTPStatusConflict                              uintptr = 409
	HTTPStatusGone                                  uintptr = 410
	HTTPStatusLengthRequired                        uintptr = 411
	HTTPStatusPrecondFailed                         uintptr = 412
	HTTPStatusRequestTooLarge                       uintptr = 413
	HTTPStatusUriTooLong                            uintptr = 414
	HTTPStatusUnsupportedMedia                      uintptr = 415
	HTTPStatusRetryWith                             uintptr = 449
	HTTPStatusServerError                           uintptr = 500
	HTTPStatusNotSupported                          uintptr = 501
	HTTPStatusBadGateway                            uintptr = 502
	HTTPStatusServiceUnavail                        uintptr = 503
	HTTPStatusGatewayTimeout                        uintptr = 504
	HTTPStatusVersionNotSup                         uintptr = 505
	HTTPStatusFirst                                 uintptr = HTTPStatusContinue
	HTTPStatusLast                                  uintptr = HTTPStatusVersionNotSup
	SecurityFlagIgnoreUnknownCa                     uintptr = 0x00000100
	SecurityFlagIgnoreCertDateInvalid               uintptr = 0x00002000
	SecurityFlagIgnoreCertCnInvalid                 uintptr = 0x00001000
	SecurityFlagIgnoreCertWrongUsage                uintptr = 0x00000200
	SecurityFlagSecure                              uintptr = 0x00000001
	SecurityFlagStrengthWeak                        uintptr = 0x10000000
	SecurityFlagStrengthMedium                      uintptr = 0x40000000
	SecurityFlagStrengthStrong                      uintptr = 0x20000000
	IcuNoEncode                                     uintptr = 0x20000000
	IcuDecode                                       uintptr = 0x10000000
	IcuNoMeta                                       uintptr = 0x08000000
	IcuEncodeSpacesOnly                             uintptr = 0x04000000
	IcuBrowserMode                                  uintptr = 0x02000000
	IcuEncodePercent                                uintptr = 0x00001000
	WinhttpQueryMimeVersion                         uintptr = 0
	WinhttpQueryContentType                         uintptr = 1
	WinhttpQueryContentTransferEncoding             uintptr = 2
	WinhttpQueryContentId                           uintptr = 3
	WinhttpQueryContentDescription                  uintptr = 4
	WinhttpQueryContentLength                       uintptr = 5
	WinhttpQueryContentLanguage                     uintptr = 6
	WinhttpQueryAllow                               uintptr = 7
	WinhttpQueryPublic                              uintptr = 8
	WinhttpQueryDate                                uintptr = 9
	WinhttpQueryExpires                             uintptr = 10
	WinhttpQueryLastModified                        uintptr = 11
	WinhttpQueryMessageId                           uintptr = 12
	WinhttpQueryUri                                 uintptr = 13
	WinhttpQueryDerivedFrom                         uintptr = 14
	WinhttpQueryCost                                uintptr = 15
	WinhttpQueryLink                                uintptr = 16
	WinhttpQueryPragma                              uintptr = 17
	WinhttpQueryVersion                             uintptr = 18
	WinhttpQueryStatusCode                          uintptr = 19
	WinhttpQueryStatusText                          uintptr = 20
	WinhttpQueryRawHeaders                          uintptr = 21
	WinhttpQueryRawHeadersCRLF                      uintptr = 22
	WinhttpQueryConnection                          uintptr = 23
	WinhttpQueryAccept                              uintptr = 24
	WinhttpQueryAcceptCharset                       uintptr = 25
	WinhttpQueryAcceptEncoding                      uintptr = 26
	WinhttpQueryAcceptLanguage                      uintptr = 27
	WinhttpQueryAuthorization                       uintptr = 28
	WinhttpQueryContentEncoding                     uintptr = 29
	WinhttpQueryForwarded                           uintptr = 30
	WinhttpQueryFrom                                uintptr = 31
	WinhttpQueryIfModifiedSince                     uintptr = 32
	WinhttpQueryLocation                            uintptr = 33
	WinhttpQueryOrigUri                             uintptr = 34
	WinhttpQueryReferer                             uintptr = 35
	WinhttpQueryRetryAfter                          uintptr = 36
	WinhttpQueryServer                              uintptr = 37
	WinhttpQueryTitle                               uintptr = 38
	WinhttpQueryUserAgent                           uintptr = 39
	WinhttpQueryWwwAuthenticate                     uintptr = 40
	WinhttpQueryProxyAuthenticate                   uintptr = 41
	WinhttpQueryAcceptRanges                        uintptr = 42
	WinhttpQuerySetCookie                           uintptr = 43
	WinhttpQueryCookie                              uintptr = 44
	WinhttpQueryRequestMethod                       uintptr = 45
	WinhttpQueryRefresh                             uintptr = 46
	WinhttpQueryContentDisposition                  uintptr = 47
	WinhttpQueryAge                                 uintptr = 48
	WinhttpQueryCacheControl                        uintptr = 49
	WinhttpQueryContentBase                         uintptr = 50
	WinhttpQueryContentLocation                     uintptr = 51
	WinhttpQueryContentMd5                          uintptr = 52
	WinhttpQueryContentRange                        uintptr = 53
	WinhttpQueryEtag                                uintptr = 54
	WinhttpQueryHost                                uintptr = 55
	WinhttpQueryIfMatch                             uintptr = 56
	WinhttpQueryIfNoneMatch                         uintptr = 57
	WinhttpQueryIfRange                             uintptr = 58
	WinhttpQueryIfUnmodifiedSince                   uintptr = 59
	WinhttpQueryMaxForwards                         uintptr = 60
	WinhttpQueryProxyAuthorization                  uintptr = 61
	WinhttpQueryRange                               uintptr = 62
	WinhttpQueryTransferEncoding                    uintptr = 63
	WinhttpQueryUpgrade                             uintptr = 64
	WinhttpQueryVary                                uintptr = 65
	WinhttpQueryVia                                 uintptr = 66
	WinhttpQueryWarning                             uintptr = 67
	WinhttpQueryExpect                              uintptr = 68
	WinhttpQueryProxyConnection                     uintptr = 69
	WinhttpQueryUnlessModifiedSince                 uintptr = 70
	WinhttpQueryProxySupport                        uintptr = 75
	WinhttpQueryAuthenticationInfo                  uintptr = 76
	WinhttpQueryPassportUrls                        uintptr = 77
	WinhttpQueryPassportConfig                      uintptr = 78
	WinhttpQueryMax                                 uintptr = 78
	WinhttpQueryCustom                              uintptr = 65535
	WinhttpQueryFlagRequestHeaders                  uintptr = 0x80000000
	WinhttpQueryFlagSystemtime                      uintptr = 0x40000000
	WinhttpQueryFlagNumber                          uintptr = 0x20000000
	WinhttpQueryFlagNumber64                        uintptr = 0x08000000
	WinhttpCallbackStatusResolvingName              uintptr = 0x00000001
	WinhttpCallbackStatusNameResolved               uintptr = 0x00000002
	WinhttpCallbackStatusConnectingToServer         uintptr = 0x00000004
	WinhttpCallbackStatusConnectedToServer          uintptr = 0x00000008
	WinhttpCallbackStatusSendingRequest             uintptr = 0x00000010
	WinhttpCallbackStatusRequestSent                uintptr = 0x00000020
	WinhttpCallbackStatusReceivingResponse          uintptr = 0x00000040
	WinhttpCallbackStatusResponseReceived           uintptr = 0x00000080
	WinhttpCallbackStatusClosingConnection          uintptr = 0x00000100
	WinhttpCallbackStatusConnectionClosed           uintptr = 0x00000200
	WinhttpCallbackStatusHandleCreated              uintptr = 0x00000400
	WinhttpCallbackStatusHandleClosing              uintptr = 0x00000800
	WinhttpCallbackStatusDetectingProxy             uintptr = 0x00001000
	WinhttpCallbackStatusRedirect                   uintptr = 0x00004000
	WinhttpCallbackStatusIntermediateResponse       uintptr = 0x00008000
	WinhttpCallbackStatusSecureFailure              uintptr = 0x00010000
	WinhttpCallbackStatusHeadersAvailable           uintptr = 0x00020000
	WinhttpCallbackStatusDataAvailable              uintptr = 0x00040000
	WinhttpCallbackStatusReadComplete               uintptr = 0x00080000
	WinhttpCallbackStatusWriteComplete              uintptr = 0x00100000
	WinhttpCallbackStatusRequestError               uintptr = 0x00200000
	WinhttpCallbackStatusSendrequestComplete        uintptr = 0x00400000
	WinhttpCallbackStatusGetproxyforurlComplete     uintptr = 0x01000000
	WinhttpCallbackStatusCloseComplete              uintptr = 0x02000000
	WinhttpCallbackStatusShutdownComplete           uintptr = 0x04000000
	WinhttpCallbackStatusSettingsWriteComplete      uintptr = 0x10000000
	WinhttpCallbackStatusSettingsReadComplete       uintptr = 0x20000000
	WinhttpCallbackFlagResolveName                  uintptr = (WinhttpCallbackStatusResolvingName | WinhttpCallbackStatusNameResolved)
	WinhttpCallbackFlagConnectToServer              uintptr = (WinhttpCallbackStatusConnectingToServer | WinhttpCallbackStatusConnectedToServer)
	WinhttpCallbackFlagSendRequest                  uintptr = (WinhttpCallbackStatusSendingRequest | WinhttpCallbackStatusRequestSent)
	WinhttpCallbackFlagReceiveResponse              uintptr = (WinhttpCallbackStatusReceivingResponse | WinhttpCallbackStatusResponseReceived)
	WinhttpCallbackFlagCloseConnection              uintptr = (WinhttpCallbackStatusClosingConnection | WinhttpCallbackStatusConnectionClosed)
	WinhttpCallbackFlagHandles                      uintptr = (WinhttpCallbackStatusHandleCreated | WinhttpCallbackStatusHandleClosing)
	WinhttpCallbackFlagDetectingProxy               uintptr = WinhttpCallbackStatusDetectingProxy
	WinhttpCallbackFlagRedirect                     uintptr = WinhttpCallbackStatusRedirect
	WinhttpCallbackFlagIntermediateResponse         uintptr = WinhttpCallbackStatusIntermediateResponse
	WinhttpCallbackFlagSecureFailure                uintptr = WinhttpCallbackStatusSecureFailure
	WinhttpCallbackFlagSendrequestComplete          uintptr = WinhttpCallbackStatusSendrequestComplete
	WinhttpCallbackFlagHeadersAvailable             uintptr = WinhttpCallbackStatusHeadersAvailable
	WinhttpCallbackFlagDataAvailable                uintptr = WinhttpCallbackStatusDataAvailable
	WinhttpCallbackFlagReadComplete                 uintptr = WinhttpCallbackStatusReadComplete
	WinhttpCallbackFlagWriteComplete                uintptr = WinhttpCallbackStatusWriteComplete
	WinhttpCallbackFlagRequestError                 uintptr = WinhttpCallbackStatusRequestError
	WinhttpCallbackFlagGetproxyforurlComplete       uintptr = WinhttpCallbackStatusGetproxyforurlComplete
	WinhttpCallbackFlagAllCompletions               uintptr = (WinhttpCallbackStatusSendrequestComplete | WinhttpCallbackStatusHeadersAvailable | WinhttpCallbackStatusDataAvailable | WinhttpCallbackStatusReadComplete | WinhttpCallbackStatusWriteComplete | WinhttpCallbackStatusRequestError | WinhttpCallbackStatusGetproxyforurlComplete)
	WinhttpCallbackFlagAllNotifications             uintptr = 0xffffffff
	ApiReceiveResponse                              uintptr = (1)
	ApiQueryDataAvailable                           uintptr = (2)
	ApiReadData                                     uintptr = (3)
	ApiWriteData                                    uintptr = (4)
	ApiSendRequest                                  uintptr = (5)
	ApiGetProxyForUrl                               uintptr = (6)
	WinhttpHandleTypeSession                        uintptr = 1
	WinhttpHandleTypeConnect                        uintptr = 2
	WinhttpHandleTypeRequest                        uintptr = 3
	WinhttpCallbackStatusFlagCertRevFailed          uintptr = 0x00000001
	WinhttpCallbackStatusFlagInvalidCert            uintptr = 0x00000002
	WinhttpCallbackStatusFlagCertRevoked            uintptr = 0x00000004
	WinhttpCallbackStatusFlagInvalidCa              uintptr = 0x00000008
	WinhttpCallbackStatusFlagCertCnInvalid          uintptr = 0x00000010
	WinhttpCallbackStatusFlagCertDateInvalid        uintptr = 0x00000020
	WinhttpCallbackStatusFlagCertWrongUsage         uintptr = 0x00000040
	WinhttpCallbackStatusFlagSecurityChannelError   uintptr = 0x80000000
	WinhttpFlagSecureProtocolSsl2                   uintptr = 0x00000008
	WinhttpFlagSecureProtocolSsl3                   uintptr = 0x00000020
	WinhttpFlagSecureProtocolTls1                   uintptr = 0x00000080
	WinhttpFlagSecureProtocolTls11                  uintptr = 0x00000200
	WinhttpFlagSecureProtocolTls12                  uintptr = 0x00000800
	WinhttpFlagSecureProtocolTls13                  uintptr = 0x00002000
	WinhttpFlagSecureProtocolAll                    uintptr = (WinhttpFlagSecureProtocolSsl2 | WinhttpFlagSecureProtocolSsl3 | WinhttpFlagSecureProtocolTls1)
	WinhttpAuthSchemeBasic                          uintptr = 0x00000001
	WinhttpAuthSchemeNtlm                           uintptr = 0x00000002
	WinhttpAuthSchemePassport                       uintptr = 0x00000004
	WinhttpAuthSchemeDigest                         uintptr = 0x00000008
	WinhttpAuthSchemeNegotiate                      uintptr = 0x00000010
	WinhttpAuthTargetServer                         uintptr = 0x00000000
	WinhttpAuthTargetProxy                          uintptr = 0x00000001
	WinhttpTimeFormatBufsize                        uintptr = 62
	WinhttpAutoDetectTypeDhcp                       uintptr = 0x00000001
	WinhttpAutoDetectTypeDnsA                       uintptr = 0x00000002
	WinhttpAutoproxyAutoDetect                      uintptr = 0x00000001
	WinhttpAutoproxyConfigUrl                       uintptr = 0x00000002
	WinhttpAutoproxyHostKeepcase                    uintptr = 0x00000004
	WinhttpAutoproxyHostLowercase                   uintptr = 0x00000008
	WinhttpAutoproxyAllowAutoconfig                 uintptr = 0x00000100
	WinhttpAutoproxyAllowStatic                     uintptr = 0x00000200
	WinhttpAutoproxyAllowCm                         uintptr = 0x00000400
	WinhttpAutoproxyRunInprocess                    uintptr = 0x00010000
	WinhttpAutoproxyRunOutprocessOnly               uintptr = 0x00020000
	WinhttpAutoproxyNoDirectaccess                  uintptr = 0x00040000
	WinhttpAutoproxyNoCacheClient                   uintptr = 0x00080000
	WinhttpAutoproxyNoCacheSvc                      uintptr = 0x00100000
	WinhttpAutoproxySortResults                     uintptr = 0x00400000
	NetworkingKeyBufsize                            uintptr = 128
	WinhttpRequestStatFlagTcpFastOpen               uintptr = 0x00000001
	WinhttpRequestStatFlagTlsSessionResumption      uintptr = 0x00000002
	WinhttpRequestStatFlagTlsFalseStart             uintptr = 0x00000004
	WinhttpRequestStatFlagProxyTlsSessionResumption uintptr = 0x00000008
	WinhttpRequestStatFlagProxyTlsFalseStart        uintptr = 0x00000010
	WinhttpRequestStatFlagFirstRequest              uintptr = 0x00000020
	WinhttpWebSocketMaxCloseReasonLength            uintptr = 123
	WinhttpWebSocketMinKeepaliveValue               uintptr = 15000
)
