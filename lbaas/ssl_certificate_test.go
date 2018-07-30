package lbaas

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

// Test the SSL Certificate lifecycle the create, get, delete an SSL Certificate
// and validate the fields are set as expected.
func TestAccServerCertificateLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	certClient, err := getSSLCertificateClient()
	assert.NoError(t, err)

	// CREATE

	createCertInput := &CreateSSLCertificateInput{
		Name:             "acc-test-ssl-cert-server1",
		Trusted:          false,
		Certificate:      "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZmRENDQTJTZ0F3SUJBZ0lRUlRvQWc1K1U0ZFhEZXltaVlBWVdCekFOQmdrcWhraUc5dzBCQVFzRkFEQXoKTVJVd0V3WURWUVFLRXd4elkzSnZjM052Y21GamJHVXhHakFZQmdOVkJBTVRFWE5qY205emMyOXlZV05zWlM1egphWFJsTUI0WERURTRNRFl6TURBeU16RXpNbG9YRFRFNU1EY3dOREEyTXpFek1sb3dVakVMTUFrR0ExVUVCaE1DClEwRXhFREFPQmdOVkJBZ1RCMDl1ZEdGeWFXOHhGVEFUQmdOVkJBb1RESE5qY205emMyOXlZV05zWlRFYU1CZ0cKQTFVRUF4TVJjMk55YjNOemIzSmhZMnhsTG5OcGRHVXdnZ0lpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElDRHdBdwpnZ0lLQW9JQ0FRQ3lUM1RadTdndE1CS3hPbXBZWjZjekRFSFVPcDVINDhqOVNnRGZUTnFQaFprR3liWnBxbC9FCmlrUnREdmxweXRxTm5yVGRhL3V1UlJibW9QWnJuUkRJNEpqT1Bybm9jUGJxcEk0U2t3c1RYdEpEbzZWTXlONUgKUUU1RFU1SDlxR0c5SWxtY2Rvc3k0dDJ2L3JNdGVvTkhLd3VVZFhkbFZzRkpGcG03cEJaaVl5M1phZzgrQjRSOApNREMwNVU1SndjMFZ3Q0hFbk1lOGhibEVsNVhJVmg2dzM4V1RRb1hESlhIWlJ1bTNRTTZTTlVaYk9MbVFMc0E4CjBhcmVmb3NqK3ZWa21qazZPYnhKU0ZGMEhnMHhNQ0Q2eUdxTzg1eG43V09obDBiYTZEVkZKU2lITjc1NnhmYTMKZzVRS0wyelZybEpaWkk4SkplNi8vcFg2QUpZQ3lyUjFQS0pVZDMxR1lHVmdBN3NvQzM1V2dsdWtqRmhxUjdwagpiNkxQbXdUM2lidFRscFNzcythU290UC81K3pBWXdTQVBCaHY2V0xlVkd5dElGaXZSczdRaVVJazVqQ1BRL1ZtCk1ScEdTMnFoTkxMM1lObW1vK2FTSzlGVzhpckplNkZLaTlJMzZ1QjR3OUF2NHA0T1FFQ3V6dlJubXhlZ1grRGEKbm9scGREc01zakFVY1gzRDRDVWVBdVhWUjczcUg2ekVmRWJZQ1FYT1ZKOE9kM3JyYzBBekVDeS9ubFJwYVlMQwpxRmlaZVExKytGeDEwaStPUXc2SmNBTUFPKzNwQjlRRDJ5NGFMZkpId3VYV0FvQUZoZFd4cVpnUGdlUE5DMHBZClpCVlg2UzBtVnJ0eXNXd2F6Z1Z0cWcxazJreFV6dXBwYTlncHhJUit4T21FODQ0eFV6MXRkd0lEQVFBQm8yMHcKYXpBVEJnTlZIU1VFRERBS0JnZ3JCZ0VGQlFjREFUQU1CZ05WSFJNQkFmOEVBakFBTUI4R0ExVWRJd1FZTUJhQQpGS3pWSnhUb1lMakl1Ujcza0VXQStIV3YwYmVnTUNVR0ExVWRFUVFlTUJ5Q0dtMTVkMlZpWVhCd0xuTmpjbTl6CmMyOXlZV05zWlM1emFYUmxNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUNBUUJtYlJpRzdYTEs0RDdvYWFaV2Q3T00KOG55aG1iSFJhY1hrZjArTG45amhjK2JCYlp6cW16WDBlQlF4WXloUWhCT1gyRjdzSDMyOU5zcmxHSEVobVNSRwplOG54S0pIV2dLT05BaURTMkU2LzdBUkVqVlJtaGJ4Z2dSZTJjUWlBUWMzR3BHbC9mTnRtT3dUdFc2MG4yYU1uCkg0WFN1Y053MnZIZy9ZMzB4bldlSVJDL1hSMkxuS1hPc2xTUkxTVThIdFk1VWRWY01MTFlUbVBJZUtDeWZ3TzEKRWhZSms1ekhyeGdvSU1NUVNHQWlrMGZmZEM1V3kvTkNiZFBLZkNhOWJjV21PWlM2YjBxMllyNm5pRkxsd1pHdQp6QklFT0ZVa1lxL3djMmttc1VYcFIrdVhxM2RPNVBqcVNVYmxjVyt0cVpvSnBQN2Y0b2tycDBiM05LUDJ1am1BCmZUcks3RUhBb0NHdWRXRGV1cEJIQnBreE9TdTlnTVp1MXQxSnJRd0M4Wlovb1I1dGkvM1B2MWVpSk9nektjaUgKaFhXYlNEL2RHQjVjblhjQWZNTVhLdngxcHAxdHJIQWdJZWYrSEU5Y2JvdkxrUWFaQXJPTzJtZ3ovbkI1Z3pJUQpySi94VnBvalhPN0c3VTJ0SCs0OU55TFl5ZE1ZWEFHRzl6UWFCRFkzZlh2M2ptQUJlKy9hWitCYkJxcWtFL3lzCnovRmhGblpYbTJTUHdCM2pmMnlRRTRhekFFZks3RmFGMVZnYXpHRGxGM1dseTN2ME8zTS9KMHhFQ3VxQTNYV1gKblVrQnNTU0poK3lkZkUzekpCWWwrSHVwTWhEazdlZlJGaW9BWTZUeDBLWUt5bXJ3Uzh1RFMvbHZ0RHlDMC9LaApQcDdMM2J4V0tXWTdUZHlpSGowYll3PT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
		PrivateKey:       "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS0FJQkFBS0NBZ0VBc2s5MDJidTRMVEFTc1RwcVdHZW5Nd3hCMURxZVIrUEkvVW9BMzB6YWo0V1pCc20yCmFhcGZ4SXBFYlE3NWFjcmFqWjYwM1d2N3JrVVc1cUQyYTUwUXlPQ1l6ajY1NkhEMjZxU09FcE1MRTE3U1E2T2wKVE1qZVIwQk9RMU9SL2FoaHZTSlpuSGFMTXVMZHIvNnpMWHFEUnlzTGxIVjNaVmJCU1JhWnU2UVdZbU10MldvUApQZ2VFZkRBd3RPVk9TY0hORmNBaHhKekh2SVc1UkplVnlGWWVzTi9GazBLRnd5VngyVWJwdDBET2tqVkdXemk1CmtDN0FQTkdxM242TEkvcjFaSm81T2ptOFNVaFJkQjROTVRBZytzaHFqdk9jWisxam9aZEcydWcxUlNVb2h6ZSsKZXNYMnQ0T1VDaTlzMWE1U1dXU1BDU1h1di82VitnQ1dBc3EwZFR5aVZIZDlSbUJsWUFPN0tBdCtWb0picEl4WQpha2U2WTIraXo1c0U5NG03VTVhVXJMUG1rcUxULytmc3dHTUVnRHdZYitsaTNsUnNyU0JZcjBiTzBJbENKT1l3CmowUDFaakVhUmt0cW9UU3k5MkRacHFQbWtpdlJWdklxeVh1aFNvdlNOK3JnZU1QUUwrS2VEa0JBcnM3MFo1c1gKb0YvZzJwNkphWFE3RExJd0ZIRjl3K0FsSGdMbDFVZTk2aCtzeEh4RzJBa0Z6bFNmRG5kNjYzTkFNeEFzdjU1VQphV21Dd3FoWW1Ya05mdmhjZGRJdmprTU9pWEFEQUR2dDZRZlVBOXN1R2kzeVI4TGwxZ0tBQllYVnNhbVlENEhqCnpRdEtXR1FWVitrdEpsYTdjckZzR3M0RmJhb05aTnBNVk03cWFXdllLY1NFZnNUcGhQT09NVk05YlhjQ0F3RUEKQVFLQ0FnQlVUakVIU1RRWldXTmRIQ3R2eFFKT3BucnhadzJ6RzhYSnpCV0Jmb3JQMVBDM1B1UGFMYzI5MVVubwo3bTJLVVhqb0FLT3ZGUUVZTWw1VGlNTTV1amRYWXFtY3loZUlDUEVWbTl2NGVFR0NWUkRCSGp4bmc0bGswc2l1CkdITXNKVktnNC83T2RWWDNKMEI5bDhVTHVhTWNJUVFHbTB0cVJJeDZqQTcvb3VOYWZWNE9MNUVwV05DUkR3L1kKVjVxZVVOMHdiWGtKeHI3QktkQ2cxN0xmMTZnSEpLWDdyRFltUUN3RitQdER3NFpucG55dllMQ0x0Uzc3RW43aApWNTlkMGFmNGV0cVg1dmhaQmJKTlhuQUtVNkVrTVdJQVdMb1lnU3JjR2ZSTVNBSDN3VXZhTXNjQ3NWcW5CYlVXCjQ3bG5obFkxSWRCbkdPSTdNSm1rYkdhQXgycHVRVkwycThJSkJRcmxWa3k1Mk9LZkZzU2pjdFEyOEUvUkJXUG4KYWh3UlQwZW5tMGtEV2J3K0JBaVcrNnQ4VEpzRFVocE03WG1aZ0hrYnFXejZaMXpQYVVmdnZLUjFHT1pyQUlyVgo3T1IyUlZHcGdkNEZDTEIvbUJoL3p3bzlPM013WFZyUUViN3daWk13Z3RwVWhXMEd0YTZaR1Y1K2lad2tkQUZuClNKZS85bVdXR3JtWER4RG9FQ1QwalVJcXNsa3hZVUVSQWo4azZxZ3hkOFVpaVFYeVREcy9abkVzQU9LMjNGYnkKUVpIakQ1WnFjdEpHczRMaEtCZ2tKNWFHcnhNcUdWRzdHT0hEYkorQ0NWaTlmQmlDdVhIOVNrZHZDamZRM21jUwpoRHlscWxFZnBzRFZ2ZmtpcHU1MDhqb3l3RHJDSDNnY21kUVZucldJTHk1VGFld3RzUUtDQVFFQTVHWXlEVXVKCjR6dHB0MExWQXJhU3Njb1dTWkltK1lKNUlIVll1aTZLVlNRQ1lvTXJTZnJqdkRRcUJJYjJnQjFiakVKUWZLMFEKc2tJa0NoK2p2ZnI1b0gycjBBc3l4VngrcmpETUhmcTV1eitoMEVBUUtITG9ueXBXUWdCWEJzRkdBdlhmL1Y0OApuWStnNWxQYVNpbmtyOGN6cGhOOFNFTENkSTUxbEdTY2hiRFd3cE9ITEgxSFhuK056bWJtZWVla2lOL2drSVhDCmVpUUppdjNPUVBWZ3dGeGZlcnp1NTdtclMrWVZERUhpcEZQV2w2SjZoejQ0OFlsM3RuQlc3bXBIMUphS05QTHkKUFRvNDlYVW41RmZBRXBsYmE3N2ZiYzV6RWQxZTQxZ2FyUVZ4endiVWtxMXFPN25qM2wyYk5aeENEano0MmRqbgpRUXNJeWUvaXBmU0pmUUtDQVFFQXg5dXo2bVVlUkVIaWw3d0ZNZU1jbGFHVThOWTF5TzdTQnJLT0pMUkxMQmNKCmlHWlEzSHRvVDFBVmtXdVNkeWlDNE81ZThYVFREWGNmQ0xPRThtK3ZBUWVNNUwwaDlPMDkxdFl3dTN5N3dUckUKWkhnbk5OM3NWVGNGNmRVdTRLdm5ZZXZoWWdaWEF6V0ZJeGpKUE1mYTc5dlZsMVN0Z0dlUTVwQ3gvcTI4YVphbQpsc2JMdVgwZG1qSXc2dUpUTkZIb2Rmd0dHWW12dkc1MTFHem00ZEdDekh6NzZOUS9wa09BQjRTRjZHZzdpWk1PCmJCbHVzcUZNbnBYOEgyc3RmWXdZQm0yZlFIems4WmhDUjZEeFVlNVFXbDB0NVcwOXRUaW0vZ3hHMnNkVTcrV0sKeGxXamtTQ0ZvdVVrejFOdG5xMzBZM1FnOUpRcWZHUVNKQzJMb1YvbEF3S0NBUUJiTmVTdklvZUNVMnU0WDl3cApKVGdZQUJnK2NUdFhVUitHTXRhb0k0WGkwbXFSWk1pWTFyU3pxREZQZFlaalMxWVFBVHViVHBIb1hqbCtRWHhtCmxoK3lVLzJWSzZPdTVXMUJxd01ZeGRQK1R6OFRwMEhNcFhiNGVUUFJUOGx4VFNYa2NNUnVybitPZkpsSTRodSsKbWxSVlRqdjJDcm9MTVgzdWhpVzJpU2RvekdJM2VpcjFQV0tPL21sbkQvamluZnM3SGd6VUtsYXI2RkJYVFZ4YwozS0V5c0xFQWx3cmhSMmg4K3ZsVTE4cm16UVJac2UwMHJVaVlUTW1kOWVjQmR6Z1FVYjRIdnkyMS9kWlpUOXdLCmVIQ2YvTlpoaE93OU1jRUtWVmxiZVFmT0tPcDJQc2dOZTJ0OVJwTVZibFJaYUhtSXJoakRCcmZ6WmJDdzFEZXoKQnFFUkFvSUJBUUMrMitkVFBzSEt1WWlsRXQ5N0pzSlRldjE3aVhYUHI1Sk81eEdycDZucUx3M2hmcVJXQ2x3dwo4ZS9HOGczclVYcTdSNmpQdVpzYnp0aUtQTFlIdC9ST2JXRjF4OUMzMENBd0hGaHBrOUxSMDBkZUV3aU9DaWo1CnNCUXJuSFNxQmtCdldRM2h5T0FycGw3QWg1a1dQRjJ1bGlmQjN4SGFBQTEyd2xQWlBSMGpVTVZDVkJLVnp4QUkKQTBxSDVSOUVaYnd6Z0R4ckF2d2FYUHFWcEhKUTBQMnlQdUZyRmRhNjl0YzdWcWx5cXFmQTEvajc3c1d5UFN1bwpmdDlKT2RjMWdDWXBiV0taK1N2Q05IK0hYQUZaRnRjUmxNNlJ2T01qUHpqcWY5cmliMTJEdzVmbGxEOGlCd2JYCjZ6QmQ5ZlJIaHlST0hjYWpDeVFQcXBsUWgxWkRCQXIxQW9JQkFIUWFFZmNVNnI1U09xc0QrVC9QSWlKZXh4T28Kb3BvMC9ZYkM2MEZIa2ZHRWRjcDVXSXdpRXNIangwd0VKbzZVWTdtOFR2aTlGcVhqTSt1Uld4RFAwZHVUaDlGVAp5QnhqN2dwRnhzbDdKTjQ2dVREb3dKVHB2UmNtZFVRRm1CbXBhY2x1dXVPSWpLNkdjMjRhRW5iRmxTemhUSGdZCnJGbUl4Q2RzamRKRDNNWndPT1kzanpyNUticXB4ZnRmaVBiblM5TzBqWFJzaTY4QW1XL29ybUpXL0hmc2NpWmsKTkR0R1BFVVgvdG9ZYVk4c1hMU1ZGU3cwY3Q3U0MxdzlUeWMxY1hwelRBdHk5NlovS3hpTnI0ZVoxb3FVVXRncAo1SEIxWGVxZVdyc0ZxbGVvQ0t3QW5jUXhKTVk3NjcwZFlER1BCL2ZVOUgrK0p4bXp1ZlNUMkY5NHFZbz0KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K",
		CertificateChain: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZNakNDQXhxZ0F3SUJBZ0lRV2hMWHVJdkhoNm9WY1NwWU5ZMmlzREFOQmdrcWhraUc5dzBCQVFzRkFEQXoKTVJVd0V3WURWUVFLRXd4elkzSnZjM052Y21GamJHVXhHakFZQmdOVkJBTVRFWE5qY205emMyOXlZV05zWlM1egphWFJsTUI0WERURTRNRFl6TURBeU16RXpNbG9YRFRFNU1EY3dOREEyTXpFek1sb3dNekVWTUJNR0ExVUVDaE1NCmMyTnliM056YjNKaFkyeGxNUm93R0FZRFZRUURFeEZ6WTNKdmMzTnZjbUZqYkdVdWMybDBaVENDQWlJd0RRWUoKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDQWdvQ2dnSUJBTEpQZE5tN3VDMHdFckU2YWxobnB6TU1RZFE2bmtmagp5UDFLQU45TTJvK0ZtUWJKdG1tcVg4U0tSRzBPK1duSzJvMmV0TjFyKzY1RkZ1YWc5bXVkRU1qZ21NNCt1ZWh3Cjl1cWtqaEtUQ3hOZTBrT2pwVXpJM2tkQVRrTlRrZjJvWWIwaVdaeDJpekxpM2EvK3N5MTZnMGNyQzVSMWQyVlcKd1VrV21idWtGbUpqTGRscUR6NEhoSHd3TUxUbFRrbkJ6UlhBSWNTY3g3eUZ1VVNYbGNoV0hyRGZ4Wk5DaGNNbApjZGxHNmJkQXpwSTFSbHM0dVpBdXdEelJxdDUraXlQNjlXU2FPVG81dkVsSVVYUWVEVEV3SVBySWFvN3puR2Z0Clk2R1hSdHJvTlVVbEtJYzN2bnJGOXJlRGxBb3ZiTld1VWxsa2p3a2w3ci8rbGZvQWxnTEt0SFU4b2xSM2ZVWmcKWldBRHV5Z0xmbGFDVzZTTVdHcEh1bU52b3MrYkJQZUp1MU9XbEt5ejVwS2kwLy9uN01CakJJQThHRy9wWXQ1VQpiSzBnV0s5R3p0Q0pRaVRtTUk5RDlXWXhHa1pMYXFFMHN2ZGcyYWFqNXBJcjBWYnlLc2w3b1VxTDBqZnE0SGpECjBDL2luZzVBUUs3TzlHZWJGNkJmNE5xZWlXbDBPd3l5TUJSeGZjUGdKUjRDNWRWSHZlb2ZyTVI4UnRnSkJjNVUKbnc1M2V1dHpRRE1RTEwrZVZHbHBnc0tvV0psNURYNzRYSFhTTDQ1RERvbHdBd0E3N2VrSDFBUGJMaG90OGtmQwo1ZFlDZ0FXRjFiR3BtQStCNDgwTFNsaGtGVmZwTFNaV3UzS3hiQnJPQlcycURXVGFURlRPNm1scjJDbkVoSDdFCjZZVHpqakZUUFcxM0FnTUJBQUdqUWpCQU1BNEdBMVVkRHdFQi93UUVBd0lDQkRBUEJnTlZIUk1CQWY4RUJUQUQKQVFIL01CMEdBMVVkRGdRV0JCU3MxU2NVNkdDNHlMa2U5NUJGZ1BoMXI5RzNvREFOQmdrcWhraUc5dzBCQVFzRgpBQU9DQWdFQXExei9LTjROZjVKNy85SjNvUUpiSmVCd0tBT1dNZDlIb2tVOEl2dDVoSXYyeTN2VW8zS2NIM1E1CnRSa0JvYWMwY2J3NlhLUlV4MkJDSFZOV3NCMHVEdEU1c3hHQ2NQc1ZlSURyV1pBK1ZiT2NEUDdKMmptcHlqV3IKaG14ZmVQVllmeTdoUUZzcWlpYUMwdjBxaHlvVXMrbmhaVXNqL3FpSVRYZlM4RGZhTUJpQ2ljQWU5dnB2dHdyRwoyTjVZM1lMS0JLUGVuQzlCVXgvNzJHaXBmRDNkYUpRVi9FbFNDQjRDVkRXaW5hVFJuWFRpRWhrNjMvMmUxeWVECmpoSjZZaWNYSS9rYytDVTk1aC8zSStQcjlPL3J4c0lhcHdQeEc0cHphUEZvZzErWm9FblJxekg1aHIwTjdESFUKZWYxdGFvZzQyVnY1bEhqSFhDYTBoNzJyKzNVSFM2MHJqdnRuQTNrclp4WXczTnEzOVJzeEhXQzgwZmx6bUhqeApBWTNkVFBKbmtMU1E4TzRlcllOdGxzbEdxRURhYmJOTUlOaU9kV05ndHZocWNwdkRUdGJ5YzlWQ25LUGhUNStOCnFOVHczUDgyTk12a0pwVzl3V1h3RVNtWkJCbVBPenAwdFdVM0tSdTNWSTBDNU5ONXZJNWdrYUphZHN4SDhPZVgKWm1OMlI1V0YrMFREYmNiVUFQUGxrQnZnelFockh0MmhlRGJmRnZRRFgzV3JoTTFIRDdLSzFKNm5nemVXVUVLUwo1cnQxbkI0bWtOKzdBYTRzL29ycUNTbC9lcmhJcjhTRkpyZGxOMEZtdlNPREFMYklFUzZweHM4RmhJSW5Oazc4CnR4QnowdG83a3g0QWYxdW5ocHlRaXU4ejFQeVRCaWV4YW54aVcvbDQraS9obXR1MzdBdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
	}

	_, err = certClient.CreateSSLCertificate(createCertInput)
	assert.NoError(t, err)

	defer destroySSLCertificate(t, certClient, createCertInput.Name)

	// FETCH

	resp, err := certClient.GetSSLCertificate(createCertInput.Name)
	assert.NoError(t, err)

	expected := &SSLCertificateInfo{
		Name:        createCertInput.Name,
		Certificate: createCertInput.Certificate,
		Trusted:     createCertInput.Trusted,
	}

	assert.Equal(t, expected.Name, resp.Name, "SSL Certificate name should match")
	assert.False(t, resp.Trusted, "SSL Certificate should not be Trusted")

}

// Test the SSL Certificate lifecycle the create, get, delete an SSL Certificate
// and validate the fields are set as expected.
// TODO Test disabled as Trusted Certificates cannot be deleted due to a know issue
// func TestAccTrustedCertificateLifeCycle(t *testing.T) {
//
// 	helper.Test(t, helper.TestCase{})
//
// 	certClient, err := getSSLCertificateClient()
// 	assert.NoError(t, err)
//
// 	// CREATE
//
// 	createCertInput := &CreateSSLCertificateInput{
// 		Name:             "acc-test-ssl-cert-trusted1",
// 		Trusted:          true,
// 		Certificate:      "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZmRENDQTJTZ0F3SUJBZ0lRUlRvQWc1K1U0ZFhEZXltaVlBWVdCekFOQmdrcWhraUc5dzBCQVFzRkFEQXoKTVJVd0V3WURWUVFLRXd4elkzSnZjM052Y21GamJHVXhHakFZQmdOVkJBTVRFWE5qY205emMyOXlZV05zWlM1egphWFJsTUI0WERURTRNRFl6TURBeU16RXpNbG9YRFRFNU1EY3dOREEyTXpFek1sb3dVakVMTUFrR0ExVUVCaE1DClEwRXhFREFPQmdOVkJBZ1RCMDl1ZEdGeWFXOHhGVEFUQmdOVkJBb1RESE5qY205emMyOXlZV05zWlRFYU1CZ0cKQTFVRUF4TVJjMk55YjNOemIzSmhZMnhsTG5OcGRHVXdnZ0lpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElDRHdBdwpnZ0lLQW9JQ0FRQ3lUM1RadTdndE1CS3hPbXBZWjZjekRFSFVPcDVINDhqOVNnRGZUTnFQaFprR3liWnBxbC9FCmlrUnREdmxweXRxTm5yVGRhL3V1UlJibW9QWnJuUkRJNEpqT1Bybm9jUGJxcEk0U2t3c1RYdEpEbzZWTXlONUgKUUU1RFU1SDlxR0c5SWxtY2Rvc3k0dDJ2L3JNdGVvTkhLd3VVZFhkbFZzRkpGcG03cEJaaVl5M1phZzgrQjRSOApNREMwNVU1SndjMFZ3Q0hFbk1lOGhibEVsNVhJVmg2dzM4V1RRb1hESlhIWlJ1bTNRTTZTTlVaYk9MbVFMc0E4CjBhcmVmb3NqK3ZWa21qazZPYnhKU0ZGMEhnMHhNQ0Q2eUdxTzg1eG43V09obDBiYTZEVkZKU2lITjc1NnhmYTMKZzVRS0wyelZybEpaWkk4SkplNi8vcFg2QUpZQ3lyUjFQS0pVZDMxR1lHVmdBN3NvQzM1V2dsdWtqRmhxUjdwagpiNkxQbXdUM2lidFRscFNzcythU290UC81K3pBWXdTQVBCaHY2V0xlVkd5dElGaXZSczdRaVVJazVqQ1BRL1ZtCk1ScEdTMnFoTkxMM1lObW1vK2FTSzlGVzhpckplNkZLaTlJMzZ1QjR3OUF2NHA0T1FFQ3V6dlJubXhlZ1grRGEKbm9scGREc01zakFVY1gzRDRDVWVBdVhWUjczcUg2ekVmRWJZQ1FYT1ZKOE9kM3JyYzBBekVDeS9ubFJwYVlMQwpxRmlaZVExKytGeDEwaStPUXc2SmNBTUFPKzNwQjlRRDJ5NGFMZkpId3VYV0FvQUZoZFd4cVpnUGdlUE5DMHBZClpCVlg2UzBtVnJ0eXNXd2F6Z1Z0cWcxazJreFV6dXBwYTlncHhJUit4T21FODQ0eFV6MXRkd0lEQVFBQm8yMHcKYXpBVEJnTlZIU1VFRERBS0JnZ3JCZ0VGQlFjREFUQU1CZ05WSFJNQkFmOEVBakFBTUI4R0ExVWRJd1FZTUJhQQpGS3pWSnhUb1lMakl1Ujcza0VXQStIV3YwYmVnTUNVR0ExVWRFUVFlTUJ5Q0dtMTVkMlZpWVhCd0xuTmpjbTl6CmMyOXlZV05zWlM1emFYUmxNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUNBUUJtYlJpRzdYTEs0RDdvYWFaV2Q3T00KOG55aG1iSFJhY1hrZjArTG45amhjK2JCYlp6cW16WDBlQlF4WXloUWhCT1gyRjdzSDMyOU5zcmxHSEVobVNSRwplOG54S0pIV2dLT05BaURTMkU2LzdBUkVqVlJtaGJ4Z2dSZTJjUWlBUWMzR3BHbC9mTnRtT3dUdFc2MG4yYU1uCkg0WFN1Y053MnZIZy9ZMzB4bldlSVJDL1hSMkxuS1hPc2xTUkxTVThIdFk1VWRWY01MTFlUbVBJZUtDeWZ3TzEKRWhZSms1ekhyeGdvSU1NUVNHQWlrMGZmZEM1V3kvTkNiZFBLZkNhOWJjV21PWlM2YjBxMllyNm5pRkxsd1pHdQp6QklFT0ZVa1lxL3djMmttc1VYcFIrdVhxM2RPNVBqcVNVYmxjVyt0cVpvSnBQN2Y0b2tycDBiM05LUDJ1am1BCmZUcks3RUhBb0NHdWRXRGV1cEJIQnBreE9TdTlnTVp1MXQxSnJRd0M4Wlovb1I1dGkvM1B2MWVpSk9nektjaUgKaFhXYlNEL2RHQjVjblhjQWZNTVhLdngxcHAxdHJIQWdJZWYrSEU5Y2JvdkxrUWFaQXJPTzJtZ3ovbkI1Z3pJUQpySi94VnBvalhPN0c3VTJ0SCs0OU55TFl5ZE1ZWEFHRzl6UWFCRFkzZlh2M2ptQUJlKy9hWitCYkJxcWtFL3lzCnovRmhGblpYbTJTUHdCM2pmMnlRRTRhekFFZks3RmFGMVZnYXpHRGxGM1dseTN2ME8zTS9KMHhFQ3VxQTNYV1gKblVrQnNTU0poK3lkZkUzekpCWWwrSHVwTWhEazdlZlJGaW9BWTZUeDBLWUt5bXJ3Uzh1RFMvbHZ0RHlDMC9LaApQcDdMM2J4V0tXWTdUZHlpSGowYll3PT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
// 		CertificateChain: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZNakNDQXhxZ0F3SUJBZ0lRV2hMWHVJdkhoNm9WY1NwWU5ZMmlzREFOQmdrcWhraUc5dzBCQVFzRkFEQXoKTVJVd0V3WURWUVFLRXd4elkzSnZjM052Y21GamJHVXhHakFZQmdOVkJBTVRFWE5qY205emMyOXlZV05zWlM1egphWFJsTUI0WERURTRNRFl6TURBeU16RXpNbG9YRFRFNU1EY3dOREEyTXpFek1sb3dNekVWTUJNR0ExVUVDaE1NCmMyTnliM056YjNKaFkyeGxNUm93R0FZRFZRUURFeEZ6WTNKdmMzTnZjbUZqYkdVdWMybDBaVENDQWlJd0RRWUoKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDQWdvQ2dnSUJBTEpQZE5tN3VDMHdFckU2YWxobnB6TU1RZFE2bmtmagp5UDFLQU45TTJvK0ZtUWJKdG1tcVg4U0tSRzBPK1duSzJvMmV0TjFyKzY1RkZ1YWc5bXVkRU1qZ21NNCt1ZWh3Cjl1cWtqaEtUQ3hOZTBrT2pwVXpJM2tkQVRrTlRrZjJvWWIwaVdaeDJpekxpM2EvK3N5MTZnMGNyQzVSMWQyVlcKd1VrV21idWtGbUpqTGRscUR6NEhoSHd3TUxUbFRrbkJ6UlhBSWNTY3g3eUZ1VVNYbGNoV0hyRGZ4Wk5DaGNNbApjZGxHNmJkQXpwSTFSbHM0dVpBdXdEelJxdDUraXlQNjlXU2FPVG81dkVsSVVYUWVEVEV3SVBySWFvN3puR2Z0Clk2R1hSdHJvTlVVbEtJYzN2bnJGOXJlRGxBb3ZiTld1VWxsa2p3a2w3ci8rbGZvQWxnTEt0SFU4b2xSM2ZVWmcKWldBRHV5Z0xmbGFDVzZTTVdHcEh1bU52b3MrYkJQZUp1MU9XbEt5ejVwS2kwLy9uN01CakJJQThHRy9wWXQ1VQpiSzBnV0s5R3p0Q0pRaVRtTUk5RDlXWXhHa1pMYXFFMHN2ZGcyYWFqNXBJcjBWYnlLc2w3b1VxTDBqZnE0SGpECjBDL2luZzVBUUs3TzlHZWJGNkJmNE5xZWlXbDBPd3l5TUJSeGZjUGdKUjRDNWRWSHZlb2ZyTVI4UnRnSkJjNVUKbnc1M2V1dHpRRE1RTEwrZVZHbHBnc0tvV0psNURYNzRYSFhTTDQ1RERvbHdBd0E3N2VrSDFBUGJMaG90OGtmQwo1ZFlDZ0FXRjFiR3BtQStCNDgwTFNsaGtGVmZwTFNaV3UzS3hiQnJPQlcycURXVGFURlRPNm1scjJDbkVoSDdFCjZZVHpqakZUUFcxM0FnTUJBQUdqUWpCQU1BNEdBMVVkRHdFQi93UUVBd0lDQkRBUEJnTlZIUk1CQWY4RUJUQUQKQVFIL01CMEdBMVVkRGdRV0JCU3MxU2NVNkdDNHlMa2U5NUJGZ1BoMXI5RzNvREFOQmdrcWhraUc5dzBCQVFzRgpBQU9DQWdFQXExei9LTjROZjVKNy85SjNvUUpiSmVCd0tBT1dNZDlIb2tVOEl2dDVoSXYyeTN2VW8zS2NIM1E1CnRSa0JvYWMwY2J3NlhLUlV4MkJDSFZOV3NCMHVEdEU1c3hHQ2NQc1ZlSURyV1pBK1ZiT2NEUDdKMmptcHlqV3IKaG14ZmVQVllmeTdoUUZzcWlpYUMwdjBxaHlvVXMrbmhaVXNqL3FpSVRYZlM4RGZhTUJpQ2ljQWU5dnB2dHdyRwoyTjVZM1lMS0JLUGVuQzlCVXgvNzJHaXBmRDNkYUpRVi9FbFNDQjRDVkRXaW5hVFJuWFRpRWhrNjMvMmUxeWVECmpoSjZZaWNYSS9rYytDVTk1aC8zSStQcjlPL3J4c0lhcHdQeEc0cHphUEZvZzErWm9FblJxekg1aHIwTjdESFUKZWYxdGFvZzQyVnY1bEhqSFhDYTBoNzJyKzNVSFM2MHJqdnRuQTNrclp4WXczTnEzOVJzeEhXQzgwZmx6bUhqeApBWTNkVFBKbmtMU1E4TzRlcllOdGxzbEdxRURhYmJOTUlOaU9kV05ndHZocWNwdkRUdGJ5YzlWQ25LUGhUNStOCnFOVHczUDgyTk12a0pwVzl3V1h3RVNtWkJCbVBPenAwdFdVM0tSdTNWSTBDNU5ONXZJNWdrYUphZHN4SDhPZVgKWm1OMlI1V0YrMFREYmNiVUFQUGxrQnZnelFockh0MmhlRGJmRnZRRFgzV3JoTTFIRDdLSzFKNm5nemVXVUVLUwo1cnQxbkI0bWtOKzdBYTRzL29ycUNTbC9lcmhJcjhTRkpyZGxOMEZtdlNPREFMYklFUzZweHM4RmhJSW5Oazc4CnR4QnowdG83a3g0QWYxdW5ocHlRaXU4ejFQeVRCaWV4YW54aVcvbDQraS9obXR1MzdBdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
// 	}
//
// 	_, err = certClient.CreateSSLCertificate(createCertInput)
// 	assert.NoError(t, err)
//
// 	defer destroySSLCertificate(t, certClient, createCertInput.Name)
//
// 	// FETCH
//
// 	resp, err := certClient.GetSSLCertificate(createCertInput.Name)
// 	assert.NoError(t, err)
//
// 	expected := &SSLCertificateInfo{
// 		Name:        createCertInput.Name,
// 		Certificate: createCertInput.Certificate,
// 		Trusted:     createCertInput.Trusted,
// 	}
//
// 	assert.Equal(t, expected.Name, resp.Name, "SSL Certificate name should match")
// 	assert.True(t, resp.Trusted, "SSL Certificate should be Trusted")
//
// }

func getSSLCertificateClient() (*SSLCertificateClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &SSLCertificateClient{}, err
	}
	return client.SSLCertificateClient(), nil
}

func destroySSLCertificate(t *testing.T, client *SSLCertificateClient, name string) {
	if _, err := client.DeleteSSLCertificate(name); err != nil {
		t.Fatal(err)
	}
}
