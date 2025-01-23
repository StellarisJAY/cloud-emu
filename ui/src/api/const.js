const Configs = {
    ApiServer: "http://localhost:8010/api/v1", // 设置后端API服务器地址
    StunServer: "stun:43.138.153.172:3478",   // 设置STUN服务器地址
    TurnServer: {
        Host: "turn:localhost:3478", // 设置TURN服务器地址
        Username: "turn_user",       // 设置TURN服务器用户名
        Password: "turn_password",   // 设置TURN服务器密码
    },

    defaultButtonLayouts: {
        "NES":[{"label": "←", "codes": ["Left"], "slot": "l4"},
            {"label": "→", "codes": ["Right"], "slot": "l6"},
            {"label": "↑", "codes": ["Up"], "slot": "l2"},
            {"label": "↓", "codes": ["Down"], "slot": "l8"},
            {"label": "A", "codes": ["A"], "slot": "r1"},
            {"label": "B", "codes": ["B"], "slot": "r7"},
            {"label": "Select", "codes": ["Select"], "slot": "r3"},
            {"label": "Start", "codes": ["Start"], "slot": "r9"},
            {"label": "↗", "codes": ["Right", "Up"], "slot": "l3"},
            {"label": "↘", "codes": ["Right", "Down"], "slot": "l9"},
            {"label": "↖", "codes": ["Left", "Up"], "slot": "l1"},
            {"label": "↙", "codes": ["Left", "Down"], "slot": "l7"},
        ],
        "GB":[{"label": "←", "codes": ["Left"], "slot": "l4"},
            {"label": "→", "codes": ["Right"], "slot": "l6"},
            {"label": "↑", "codes": ["Up"], "slot": "l2"},
            {"label": "↓", "codes": ["Down"], "slot": "l8"},
            {"label": "A", "codes": ["A"], "slot": "r1"},
            {"label": "B", "codes": ["B"], "slot": "r7"},
            {"label": "Select", "codes": ["Select"], "slot": "r3"},
            {"label": "Start", "codes": ["Start"], "slot": "r9"},
            {"label": "↗", "codes": ["Right", "Up"], "slot": "l3"},
            {"label": "↘", "codes": ["Right", "Down"], "slot": "l9"},
            {"label": "↖", "codes": ["Left", "Up"], "slot": "l1"},
            {"label": "↙", "codes": ["Left", "Down"], "slot": "l7"},
        ],
        "GBA":[{"label": "←", "codes": ["Left"], "slot": "l4"},
            {"label": "→", "codes": ["Right"], "slot": "l6"},
            {"label": "↑", "codes": ["Up"], "slot": "l2"},
            {"label": "↓", "codes": ["Down"], "slot": "l8"},
            {"label": "A", "codes": ["A"], "slot": "r1"},
            {"label": "B", "codes": ["B"], "slot": "r7"},
            {"label": "Select", "codes": ["Select"], "slot": "r3"},
            {"label": "Start", "codes": ["Start"], "slot": "r9"},
            {"label": "↗", "codes": ["Right", "Up"], "slot": "l3"},
            {"label": "↘", "codes": ["Right", "Down"], "slot": "l9"},
            {"label": "↖", "codes": ["Left", "Up"], "slot": "l1"},
            {"label": "↙", "codes": ["Left", "Down"], "slot": "l7"},
        ],
        "GBC":[{"label": "←", "codes": ["Left"], "slot": "l4"},
            {"label": "→", "codes": ["Right"], "slot": "l6"},
            {"label": "↑", "codes": ["Up"], "slot": "l2"},
            {"label": "↓", "codes": ["Down"], "slot": "l8"},
            {"label": "A", "codes": ["A"], "slot": "r1"},
            {"label": "B", "codes": ["B"], "slot": "r7"},
            {"label": "Select", "codes": ["Select"], "slot": "r3"},
            {"label": "Start", "codes": ["Start"], "slot": "r9"},
            {"label": "↗", "codes": ["Right", "Up"], "slot": "l3"},
            {"label": "↘", "codes": ["Right", "Down"], "slot": "l9"},
            {"label": "↖", "codes": ["Left", "Up"], "slot": "l1"},
            {"label": "↙", "codes": ["Left", "Down"], "slot": "l7"},
        ],
    },
    defaultKeyboardBindings: {
        "NES": [
            {keyCode: "KeyA", button: "Left"},
            {keyCode: "KeyD", button: "Right"},
            {keyCode: "KeyW", button: "Up"},
            {keyCode: "KeyS", button: "Down"},
            {keyCode: "Space", button: "A"},
            {keyCode: "KeyJ", button: "B"},
            {keyCode: "Tab", button: "Select"},
            {keyCode: "Enter", button: "Start"},
        ],
        "GB": [
            {keyCode: "KeyA", button: "Left"},
            {keyCode: "KeyD", button: "Right"},
            {keyCode: "KeyW", button: "Up"},
            {keyCode: "KeyS", button: "Down"},
            {keyCode: "Space", button: "A"},
            {keyCode: "KeyJ", button: "B"},
            {keyCode: "Tab", button: "Select"},
            {keyCode: "Enter", button: "Start"},
        ],
        "GBC": [
            {keyCode: "KeyA", button: "Left"},
            {keyCode: "KeyD", button: "Right"},
            {keyCode: "KeyW", button: "Up"},
            {keyCode: "KeyS", button: "Down"},
            {keyCode: "Space", button: "A"},
            {keyCode: "KeyJ", button: "B"},
            {keyCode: "Tab", button: "Select"},
            {keyCode: "Enter", button: "Start"},
        ],
        "GBA": [
            {keyCode: "KeyA", button: "Left"},
            {keyCode: "KeyD", button: "Right"},
            {keyCode: "KeyW", button: "Up"},
            {keyCode: "KeyS", button: "Down"},
            {keyCode: "Space", button: "A"},
            {keyCode: "KeyJ", button: "B"},
            {keyCode: "Tab", button: "Select"},
            {keyCode: "Enter", button: "Start"},
        ]
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