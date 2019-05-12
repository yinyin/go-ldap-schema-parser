package ldapschemaparser

import (
	"testing"
)

const sampleMatchRuleUse1 = "( 1.2.840.113556.1.4.803 NAME 'integerBitAndMatch' APPLIES ( " +
	"supportedLDAPVersion $ entryTtl $ uidNumber $ gidNumber $ olcConcurrency $ ol" +
	"cConnMaxPending $ olcConnMaxPendingAuth $ olcIdleTimeout $ olcIndexSubstrIfMi" +
	"nLen $ olcIndexSubstrIfMaxLen $ olcIndexSubstrAnyLen $ olcIndexSubstrAnyStep" +
	" $ olcIndexIntLen $ olcListenerThreads $ olcLocalSSF $ olcMaxDerefDepth $ olcR" +
	"eplicationInterval $ olcSockbufMaxIncoming $ olcSockbufMaxIncomingAuth $ olcT" +
	"hreads $ olcToolThreads $ olcWriteTimeout $ olcDbMaxReaders $ olcDbMaxSize $ " +
	"olcDbRtxnSize $ olcDbSearchStack $ mailPreferenceOption $ shadowLastChange $ " +
	"shadowMin $ shadowMax $ shadowWarning $ shadowInactive $ shadowExpire $ shado" +
	"wFlag $ ipServicePort $ ipProtocolNumber $ oncRpcNumber $ sambaPwdLastSet $ s" +
	"ambaPwdCanChange $ sambaPwdMustChange $ sambaLogonTime $ sambaLogoffTime $ sa" +
	"mbaKickoffTime $ sambaBadPasswordCount $ sambaBadPasswordTime $ sambaGroupTyp" +
	"e $ sambaNextUserRid $ sambaNextGroupRid $ sambaNextRid $ sambaAlgorithmicRid" +
	"Base $ sambaIntegerOption $ sambaMinPwdLength $ sambaPwdHistoryLength $ samba" +
	"LogonToChgPwd $ sambaMaxPwdAge $ sambaMinPwdAge $ sambaLockoutDuration $ samb" +
	"aLockoutObservationWindow $ sambaLockoutThreshold $ sambaForceLogoff $ sambaR" +
	"efuseMachinePwdChange $ sambaTrustType $ sambaTrustAttributes $ sambaTrustDir" +
	"ection $ sambaTrustPosixOffset $ sambaSupportedEncryptionTypes ) )"

const sampleMatchRuleUse2 = "( 2.5.13.30 NAME 'objectIdentifierFirstComponentMatch' APPLIE" +
	"S ( supportedControl $ supportedExtension $ supportedFeatures $ ldapSyntaxes" +
	" $ supportedApplicationContext ) )"

func TestMatchRuleSchemaUse_1(t *testing.T) {
	s, err := ParseMatchingRuleUseSchema(sampleMatchRuleUse1)
	if nil != err {
		t.Fatalf("failed on parsing Match Rule Use sample 1: %v", err)
	}
	v := s.String()
	if v != sampleMatchRuleUse1 {
		t.Errorf("expecting %v but have %v", sampleMatchRuleUse1, v)
	}
}

func TestMatchRuleSchemaUse_2(t *testing.T) {
	s, err := ParseMatchingRuleUseSchema(sampleMatchRuleUse2)
	if nil != err {
		t.Fatalf("failed on parsing Match Rule Use sample 2: %v", err)
	}
	v := s.String()
	if v != sampleMatchRuleUse2 {
		t.Errorf("expecting %v but have %v", sampleMatchRuleUse2, v)
	}
	if len(s.AppliesTo) != 5 {
		t.Errorf("expecting 5 APPLYTO items but have %v (%v)", len(s.AppliesTo), s.AppliesTo)
	}
}
