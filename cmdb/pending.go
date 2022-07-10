package cmdb

import "github.com/philcantcode/localmapper/system"

func (entity *Entity) ApprovePending() {
	entity.SysTags = append(entity.SysTags, EntityTag{Label: "Verified", DataType: system.DataType_BOOL, Values: []string{"1"}})

	entity.InsertInventory()
	entity.DELETE_ENTRY_Pending()
}
