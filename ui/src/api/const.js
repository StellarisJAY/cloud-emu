const Configs = {
    ApiServer: "http://localhost:8010/api/v1", // 设置后端API服务器地址
    StunServer: "stun:43.138.153.172:3478",   // 设置STUN服务器地址
    TurnServer: {
        Host: "turn:localhost:3478", // 设置TURN服务器地址
        Username: "turn_user",       // 设置TURN服务器用户名
        Password: "turn_password",   // 设置TURN服务器密码
    },

    defaultButtonLayouts: {
        "NES":[{"label": "←", "code": "Left", "slot": "l4"},
            {"label": "→", "code": "Right", "slot": "l6"},
            {"label": "↑", "code": "Up", "slot": "l2"},
            {"label": "↓", "code": "Down", "slot": "l8"},
            {"label": "A", "code": "A", "slot": "r1"},
            {"label": "B", "code": "B", "slot": "r7"},
            {"label": "Select", "code": "Select", "slot": "r3"},
            {"label": "Start", "code": "Start", "slot": "r9"}],
        "CHIP8": [
            {"label": "1", "code": "1", "slot": "l1"},
            {"label": "2", "code": "2", "slot": "l2"},
            {"label": "3", "code": "3", "slot": "l3"},
            {"label": "4", "code": "4", "slot": "l4"},
            {"label": "5", "code": "5", "slot": "l5"},
            {"label": "6", "code": "6", "slot": "l6"},
            {"label": "7", "code": "7", "slot": "l7"},
            {"label": "8", "code": "8", "slot": "l8"},
            {"label": "9", "code": "9", "slot": "l9"},
            {"label": "A", "code": "A", "slot": "r7"},
            {"label": "B", "code": "B", "slot": "r8"},
            {"label": "C", "code": "C", "slot": "r3"},
            {"label": "D", "code": "D", "slot": "r6"},
            {"label": "E", "code": "E", "slot": "r9"},
            {"label": "F", "code": "F", "slot": "r5"},
            {"label": "0", "code": "0", "slot": "r4"},
        ],
        "GB": [{"label": "←", "code": "Left", "slot": "l4"},
            {"label": "→", "code": "Right", "slot": "l6"},
            {"label": "↑", "code": "Up", "slot": "l2"},
            {"label": "↓", "code": "Down", "slot": "l8"},
            {"label": "A", "code": "A", "slot": "r1"},
            {"label": "B", "code": "B", "slot": "r7"},
            {"label": "Select", "code": "Select", "slot": "r3"},
            {"label": "Start", "code": "Start", "slot": "r9"}],
        "GBA": [{"label": "←", "code": "Left", "slot": "l4"},
            {"label": "→", "code": "Right", "slot": "l6"},
            {"label": "↑", "code": "Up", "slot": "l2"},
            {"label": "↓", "code": "Down", "slot": "l8"},
            {"label": "A", "code": "A", "slot": "r1"},
            {"label": "B", "code": "B", "slot": "r4"},
            {"label": "R", "code": "R", "slot": "r3"},
            {"label": "L", "code": "L", "slot": "r6"},
            {"label": "Select", "code": "Select", "slot": "r7"},
            {"label": "Start", "code": "Start", "slot": "r9"}],
    },

    enums: {
        roomJoinTypeEnum: [
            {"id": 1, "name": "公开"},
            {"id": 2, "name": "私有"},
        ],
        userStatusEnum: [
            {"id": 1, "name": "可用"},
            {"id": 2, "name": "未激活"},
            {"id": 3, "name": "封禁"}
        ],
        roomMemberRoleEnum: [
            {"id": 1, "name": "房主"},
            {"id": 2, "name": "玩家"},
            {"id": 3, "name": "观战"}
        ]
    },

    getEnum: function(enumName) {
        return this.enums[enumName];
    },
    getEnumOptions: function(enumName) {
        return this.enums[enumName].map((item) => {
            return {
                label: item.name,
                value: item.id
            }
        });
    },
    getEnumName: function(enumName, id) {
        return this.enums[enumName].find((item) => item.id === id).name;
    },
    getEnumOptionsWithAll: function(enumName) {
        let res = this.enums[enumName].map((item) => {
            return {
                label: item.name,
                value: item.id
            }
        });
        res.splice(0,0, {
            "label": "全部",
            "value": 0
        });
        return res;
    }
};

export default Configs;