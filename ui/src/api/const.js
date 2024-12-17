const Configs = {
    ApiServer: "http://localhost:8010/api/v1", // 设置后端API服务器地址
    StunServer: "stun:localhost:3478",   // 设置STUN服务器地址
    TurnServer: {
        Host: "turn:localhost:3478", // 设置TURN服务器地址
        Username: "turn_user",       // 设置TURN服务器用户名
        Password: "turn_password",   // 设置TURN服务器密码
    },
    defaultKeyboardSetting: {
        "id": "0",
        "name": "默认设置",
        "bindings": [
            {
                "emulatorKey": "Left",
                "emulatorKeyTranslated": "Left",
                "buttons": ["KeyA"],
            },
            {
                "emulatorKey": "Right",
                "emulatorKeyTranslated": "Right",
                "buttons": ["KeyD"],
            },
            {
                "emulatorKey": "Up",
                "emulatorKeyTranslated": "Up",
                "buttons": ["KeyW"],
            },
            {
                "emulatorKey": "Down",
                "emulatorKeyTranslated": "Down",
                "buttons": ["KeyS"],
            },
            {
                "emulatorKey": "A",
                "emulatorKeyTranslated": "A",
                "buttons": ["Space"],
            },
            {
                "emulatorKey": "B",
                "emulatorKeyTranslated": "B",
                "buttons": ["KeyJ"],
            },
            {
                "emulatorKey": "Start",
                "emulatorKeyTranslated": "Start",
                "buttons": ["Enter"],
            },
            {
                "emulatorKey": "Select",
                "emulatorKeyTranslated": "Select",
                "buttons": ["Tab"],
            },
        ]
    },
    
    enums: {
        roomJoinTypeEnum: [
            {"id": 1, "name": "公开"},
            {"id": 2, "name": "私有"},
            {"id": 3, "name": "仅邀请"}
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