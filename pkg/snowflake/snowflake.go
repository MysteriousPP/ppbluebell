package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	//sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
	node          *snowflake.Node
)

// func getMachineID() (uint16, error) {
// 	return sonyMachineID, nil
// }

// // 需传入当前的机器ID
// func Init(machineId uint16) (err error) {
// 	sonyMachineID = machineId
// 	t, _ := time.Parse("2006-01-02", "2020-01-01")
// 	settings := sonyflake.Settings{
// 		StartTime: t,
// 		MachineID: getMachineID,
// 	}
// 	sonyFlake = sonyflake.NewSonyflake(settings)
// 	return
// }

// // GetID 返回生成的id值
// func GetID() (id uint64, err error) {
// 	if sonyFlake == nil {
// 		err = fmt.Errorf("snoy flake not inited")
// 		return
// 	}

//		id, err = sonyFlake.NextID()
//		return
//	}
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}

func main() {
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("init failed, err:%v\n", err)
		return
	}
	id := GenID()
	fmt.Println(id)
}
