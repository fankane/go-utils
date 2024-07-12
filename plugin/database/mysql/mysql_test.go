package mysql

import (
	"fmt"
	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/plugin"
	"testing"
	"time"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if DB == nil {
		fmt.Println("db is nil")
		return
	}
	//rows, err := DB.Query("select name from big_data_test")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//for rows.Next() {
	//	var temp string
	//	if err = rows.Scan(&temp); err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println("temp:", temp)
	//}
	funcList := make([]func() error, 0)
	for i := 0; i < 1; i++ {
		funcList = append(funcList, func() error {
			return batchInsert(12)
		})
	}

	if err := goroutine.Exec(funcList); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("success")
}

func batchInsert(cnt int) error {
	// 开始事务
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//fmt.Println("tx success")

	// 准备SQL语句
	stmt, err := tx.Prepare("INSERT INTO big_data_test(`name`, `age`, `address`, `desc`, big_desc) VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 执行批量插入
	for i := 0; i < cnt; i++ {
		_, err = stmt.Exec(fmt.Sprintf("hello%d", i), i+1, "湖北省武汉市江夏区", "我是个描述", tt)
		if err != nil {
			return fmt.Errorf("exec err:%s", err)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("tx commit ", time.Now())
	return nil
}

var tt = `2023年4月13日，西安市新城区人民法院对服毒自杀的西安大三女孩程橙的生命权纠纷案作出判决。
程橙父母认为任某某不愿退钱并在明知程橙不能受刺激的情况下，多次以嘲讽、辱骂及喝农药才退款等言语刺激程橙，致其身心俱疲、万念俱灰而含恨自杀，起诉房东任某某要求赔偿40万元。庭审播放程橙生前求助情感电台与警方的录音时，她的母亲刘女士多次捂嘴啜泣。
开庭时70多岁的男房东任某某未露面，他的律师认为程橙患有双相情感障碍，任某某的言行与程橙的死亡没有因果关系，程橙去世后任某某承受了内心的煎熬及他人的责难，任某某才是真正受害者，程橙父母的不负责任应对程橙死亡负有更多责任。
双相情感障碍，是一种既有躁狂症发作，又有抑郁症发作的常见精神障碍，首次发病可见于任何年龄。当躁狂发作时，患者有情感高涨、言语活动增多、精力充沛等表现；而当抑郁发作时，患者又常表现出情绪低落、愉快感丧失、言语活动减少、疲劳迟钝等症状。
法院认为，任某某与程橙在协商退还租金及押金的过程中，对程橙进行辱骂和挖苦，给程橙的心理和精神造成极大困扰，其言行违反公序良俗。虽然任某某对程橙的自杀行为无法预见，但在程橙反复提醒不要再对其进行刺激的情况下，任某某不但没有及时注意、理性沟通，还在明知程橙解除租赁合同的原因系其与男友分手的情况下出言侮辱，其行为具有一定的过错，应当承担相应的民事赔偿责任。程橙作为成年人，对生活中遇到的纠纷，亦应及时调整心态，理性面对。
法院判决：国无德不兴，人无德不立，判任某某赔偿医疗费、丧葬费、死亡赔偿金共计174643元、精神损害抚慰金20000元，合计194643元。
近日，西安市中级人民法院经过二审，法官认为，房东辱骂房客程某的行为，给程某精神上造成了一定的伤害，存在过错，其行为也与社会主义核心价值观诚信、友善的价值准则相悖。虽程某系服毒自杀，但房东任某某的行为与程某自杀的后果存在一定的因果关系，故判决“驳回上诉，维持原判”。`
