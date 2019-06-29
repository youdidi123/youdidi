+function($){

  $.rawCitiesData = [
    {
      "name":"油田单位",
      "code":"1000000000000",
      "sub":[
        {
          "name":"采油一厂",
          "code":"1000100000000",
          "sub":[
            {
              "name":"厂部",
              "code":"1000101610824"
            },
          ]
        },
        {
          "name":"采油二厂",
          "code":"1000200000000",
          "sub":[
            {
              "name":"厂部",
              "code":"1000201621021",
            },
            {
              "name":"五蛟作业区",
              "code":"1000201621023",
            },
          ]
        },
      ]
    },
    {
      "name":"油田基地",
      "code":"2000000000000",
      "sub":[
        {
          "name":"西安片区",
          "code":"2000100000000",
          "sub":[
            {
              "name":"西安基地",
              "code":"2000101610112"
            },
            {
              "name":"龙凤园",
              "code":"2000102610126"
            },
            {
              "name":"泾渭园",
              "code":"2000103610126"
            },
          ]
        },
        {
          "name":"银川片区",
          "code":"2000200000000",
          "sub":[
            {
              "name":"燕鸽湖",
              "code":"2000201640104",
            },
          ]
        },
      ]
    },
    {
      "name": "陕西省",
      "code": "610000",
      "sub": [{
        "name": "西安市",
        "code": "610100",
        "sub": [{
          "name": "市辖区",
          "code": "610101"
        },
          {
            "name": "新城区",
            "code": "610102"
          },
          {
            "name": "碑林区",
            "code": "610103"
          },
          {
            "name": "莲湖区",
            "code": "610104"
          },
          {
            "name": "灞桥区",
            "code": "610111"
          },
          {
            "name": "未央区",
            "code": "610112"
          },
          {
            "name": "雁塔区",
            "code": "610113"
          },
          {
            "name": "阎良区",
            "code": "610114"
          },
          {
            "name": "临潼区",
            "code": "610115"
          },
          {
            "name": "长安区",
            "code": "610116"
          },
          {
            "name": "高陵区",
            "code": "610117"
          },
          {
            "name": "蓝田县",
            "code": "610122"
          },
          {
            "name": "周至县",
            "code": "610124"
          },
          {
            "name": "户县",
            "code": "610125"
          }
        ]
      },
        {
          "name": "铜川市",
          "code": "610200",
          "sub": [{
            "name": "市辖区",
            "code": "610201"
          },
            {
              "name": "王益区",
              "code": "610202"
            },
            {
              "name": "印台区",
              "code": "610203"
            },
            {
              "name": "耀州区",
              "code": "610204"
            },
            {
              "name": "宜君县",
              "code": "610222"
            }
          ]
        },
        {
          "name": "宝鸡市",
          "code": "610300",
          "sub": [{
            "name": "市辖区",
            "code": "610301"
          },
            {
              "name": "渭滨区",
              "code": "610302"
            },
            {
              "name": "金台区",
              "code": "610303"
            },
            {
              "name": "陈仓区",
              "code": "610304"
            },
            {
              "name": "凤翔县",
              "code": "610322"
            },
            {
              "name": "岐山县",
              "code": "610323"
            },
            {
              "name": "扶风县",
              "code": "610324"
            },
            {
              "name": "眉县",
              "code": "610326"
            },
            {
              "name": "陇县",
              "code": "610327"
            },
            {
              "name": "千阳县",
              "code": "610328"
            },
            {
              "name": "麟游县",
              "code": "610329"
            },
            {
              "name": "凤县",
              "code": "610330"
            },
            {
              "name": "太白县",
              "code": "610331"
            }
          ]
        },
        {
          "name": "咸阳市",
          "code": "610400",
          "sub": [{
            "name": "市辖区",
            "code": "610401"
          },
            {
              "name": "秦都区",
              "code": "610402"
            },
            {
              "name": "杨陵区",
              "code": "610403"
            },
            {
              "name": "渭城区",
              "code": "610404"
            },
            {
              "name": "三原县",
              "code": "610422"
            },
            {
              "name": "泾阳县",
              "code": "610423"
            },
            {
              "name": "乾县",
              "code": "610424"
            },
            {
              "name": "礼泉县",
              "code": "610425"
            },
            {
              "name": "永寿县",
              "code": "610426"
            },
            {
              "name": "彬县",
              "code": "610427"
            },
            {
              "name": "长武县",
              "code": "610428"
            },
            {
              "name": "旬邑县",
              "code": "610429"
            },
            {
              "name": "淳化县",
              "code": "610430"
            },
            {
              "name": "武功县",
              "code": "610431"
            },
            {
              "name": "兴平市",
              "code": "610481"
            }
          ]
        },
        {
          "name": "渭南市",
          "code": "610500",
          "sub": [{
            "name": "市辖区",
            "code": "610501"
          },
            {
              "name": "临渭区",
              "code": "610502"
            },
            {
              "name": "华县",
              "code": "610521"
            },
            {
              "name": "潼关县",
              "code": "610522"
            },
            {
              "name": "大荔县",
              "code": "610523"
            },
            {
              "name": "合阳县",
              "code": "610524"
            },
            {
              "name": "澄城县",
              "code": "610525"
            },
            {
              "name": "蒲城县",
              "code": "610526"
            },
            {
              "name": "白水县",
              "code": "610527"
            },
            {
              "name": "富平县",
              "code": "610528"
            },
            {
              "name": "韩城市",
              "code": "610581"
            },
            {
              "name": "华阴市",
              "code": "610582"
            }
          ]
        },
        {
          "name": "延安市",
          "code": "610600",
          "sub": [{
            "name": "市辖区",
            "code": "610601"
          },
            {
              "name": "宝塔区",
              "code": "610602"
            },
            {
              "name": "延长县",
              "code": "610621"
            },
            {
              "name": "延川县",
              "code": "610622"
            },
            {
              "name": "子长县",
              "code": "610623"
            },
            {
              "name": "安塞县",
              "code": "610624"
            },
            {
              "name": "志丹县",
              "code": "610625"
            },
            {
              "name": "吴起县",
              "code": "610626"
            },
            {
              "name": "甘泉县",
              "code": "610627"
            },
            {
              "name": "富县",
              "code": "610628"
            },
            {
              "name": "洛川县",
              "code": "610629"
            },
            {
              "name": "宜川县",
              "code": "610630"
            },
            {
              "name": "黄龙县",
              "code": "610631"
            },
            {
              "name": "黄陵县",
              "code": "610632"
            }
          ]
        },
        {
          "name": "汉中市",
          "code": "610700",
          "sub": [{
            "name": "市辖区",
            "code": "610701"
          },
            {
              "name": "汉台区",
              "code": "610702"
            },
            {
              "name": "南郑县",
              "code": "610721"
            },
            {
              "name": "城固县",
              "code": "610722"
            },
            {
              "name": "洋县",
              "code": "610723"
            },
            {
              "name": "西乡县",
              "code": "610724"
            },
            {
              "name": "勉县",
              "code": "610725"
            },
            {
              "name": "宁强县",
              "code": "610726"
            },
            {
              "name": "略阳县",
              "code": "610727"
            },
            {
              "name": "镇巴县",
              "code": "610728"
            },
            {
              "name": "留坝县",
              "code": "610729"
            },
            {
              "name": "佛坪县",
              "code": "610730"
            }
          ]
        },
        {
          "name": "榆林市",
          "code": "610800",
          "sub": [{
            "name": "市辖区",
            "code": "610801"
          },
            {
              "name": "榆阳区",
              "code": "610802"
            },
            {
              "name": "神木县",
              "code": "610821"
            },
            {
              "name": "府谷县",
              "code": "610822"
            },
            {
              "name": "横山县",
              "code": "610823"
            },
            {
              "name": "靖边县",
              "code": "610824"
            },
            {
              "name": "定边县",
              "code": "610825"
            },
            {
              "name": "绥德县",
              "code": "610826"
            },
            {
              "name": "米脂县",
              "code": "610827"
            },
            {
              "name": "佳县",
              "code": "610828"
            },
            {
              "name": "吴堡县",
              "code": "610829"
            },
            {
              "name": "清涧县",
              "code": "610830"
            },
            {
              "name": "子洲县",
              "code": "610831"
            }
          ]
        },
        {
          "name": "安康市",
          "code": "610900",
          "sub": [{
            "name": "市辖区",
            "code": "610901"
          },
            {
              "name": "汉阴县",
              "code": "610921"
            },
            {
              "name": "石泉县",
              "code": "610922"
            },
            {
              "name": "宁陕县",
              "code": "610923"
            },
            {
              "name": "紫阳县",
              "code": "610924"
            },
            {
              "name": "岚皋县",
              "code": "610925"
            },
            {
              "name": "平利县",
              "code": "610926"
            },
            {
              "name": "镇坪县",
              "code": "610927"
            },
            {
              "name": "旬阳县",
              "code": "610928"
            },
            {
              "name": "白河县",
              "code": "610929"
            }
          ]
        },
        {
          "name": "商洛市",
          "code": "611000",
          "sub": [{
            "name": "市辖区",
            "code": "611001"
          },
            {
              "name": "商州区",
              "code": "611002"
            },
            {
              "name": "洛南县",
              "code": "611021"
            },
            {
              "name": "丹凤县",
              "code": "611022"
            },
            {
              "name": "商南县",
              "code": "611023"
            },
            {
              "name": "山阳县",
              "code": "611024"
            },
            {
              "name": "镇安县",
              "code": "611025"
            },
            {
              "name": "柞水县",
              "code": "611026"
            }
          ]
        }
      ]
    },
    {
      "name": "甘肃省",
      "code": "620000",
      "sub": [{
        "name": "兰州市",
        "code": "620100",
        "sub": [{
          "name": "市辖区",
          "code": "620101"
        },
          {
            "name": "城关区",
            "code": "620102"
          },
          {
            "name": "七里河区",
            "code": "620103"
          },
          {
            "name": "西固区",
            "code": "620104"
          },
          {
            "name": "安宁区",
            "code": "620105"
          },
          {
            "name": "红古区",
            "code": "620111"
          },
          {
            "name": "永登县",
            "code": "620121"
          },
          {
            "name": "皋兰县",
            "code": "620122"
          },
          {
            "name": "榆中县",
            "code": "620123"
          }
        ]
      },
        {
          "name": "嘉峪关市",
          "code": "620200",
          "sub": [{
            "name": "市辖区",
            "code": "620201"
          }]
        },
        {
          "name": "金昌市",
          "code": "620300",
          "sub": [{
            "name": "市辖区",
            "code": "620301"
          },
            {
              "name": "金川区",
              "code": "620302"
            },
            {
              "name": "永昌县",
              "code": "620321"
            }
          ]
        },
        {
          "name": "白银市",
          "code": "620400",
          "sub": [{
            "name": "市辖区",
            "code": "620401"
          },
            {
              "name": "白银区",
              "code": "620402"
            },
            {
              "name": "平川区",
              "code": "620403"
            },
            {
              "name": "靖远县",
              "code": "620421"
            },
            {
              "name": "会宁县",
              "code": "620422"
            },
            {
              "name": "景泰县",
              "code": "620423"
            }
          ]
        },
        {
          "name": "天水市",
          "code": "620500",
          "sub": [{
            "name": "市辖区",
            "code": "620501"
          },
            {
              "name": "秦州区",
              "code": "620502"
            },
            {
              "name": "麦积区",
              "code": "620503"
            },
            {
              "name": "清水县",
              "code": "620521"
            },
            {
              "name": "秦安县",
              "code": "620522"
            },
            {
              "name": "甘谷县",
              "code": "620523"
            },
            {
              "name": "武山县",
              "code": "620524"
            },
            {
              "name": "张家川回族自治县",
              "code": "620525"
            }
          ]
        },
        {
          "name": "武威市",
          "code": "620600",
          "sub": [{
            "name": "市辖区",
            "code": "620601"
          },
            {
              "name": "凉州区",
              "code": "620602"
            },
            {
              "name": "民勤县",
              "code": "620621"
            },
            {
              "name": "古浪县",
              "code": "620622"
            },
            {
              "name": "天祝藏族自治县",
              "code": "620623"
            }
          ]
        },
        {
          "name": "张掖市",
          "code": "620700",
          "sub": [{
            "name": "市辖区",
            "code": "620701"
          },
            {
              "name": "甘州区",
              "code": "620702"
            },
            {
              "name": "肃南裕固族自治县",
              "code": "620721"
            },
            {
              "name": "民乐县",
              "code": "620722"
            },
            {
              "name": "临泽县",
              "code": "620723"
            },
            {
              "name": "高台县",
              "code": "620724"
            },
            {
              "name": "山丹县",
              "code": "620725"
            }
          ]
        },
        {
          "name": "平凉市",
          "code": "620800",
          "sub": [{
            "name": "市辖区",
            "code": "620801"
          },
            {
              "name": "崆峒区",
              "code": "620802"
            },
            {
              "name": "泾川县",
              "code": "620821"
            },
            {
              "name": "灵台县",
              "code": "620822"
            },
            {
              "name": "崇信县",
              "code": "620823"
            },
            {
              "name": "华亭县",
              "code": "620824"
            },
            {
              "name": "庄浪县",
              "code": "620825"
            },
            {
              "name": "静宁县",
              "code": "620826"
            }
          ]
        },
        {
          "name": "酒泉市",
          "code": "620900",
          "sub": [{
            "name": "市辖区",
            "code": "620901"
          },
            {
              "name": "肃州区",
              "code": "620902"
            },
            {
              "name": "金塔县",
              "code": "620921"
            },
            {
              "name": "瓜州县",
              "code": "620922"
            },
            {
              "name": "肃北蒙古族自治县",
              "code": "620923"
            },
            {
              "name": "阿克塞哈萨克族自治县",
              "code": "620924"
            },
            {
              "name": "玉门市",
              "code": "620981"
            },
            {
              "name": "敦煌市",
              "code": "620982"
            }
          ]
        },
        {
          "name": "庆阳市",
          "code": "621000",
          "sub": [{
            "name": "市辖区",
            "code": "621001"
          },
            {
              "name": "西峰区",
              "code": "621002"
            },
            {
              "name": "庆城县",
              "code": "621021"
            },
            {
              "name": "环县",
              "code": "621022"
            },
            {
              "name": "华池县",
              "code": "621023"
            },
            {
              "name": "合水县",
              "code": "621024"
            },
            {
              "name": "正宁县",
              "code": "621025"
            },
            {
              "name": "宁县",
              "code": "621026"
            },
            {
              "name": "镇原县",
              "code": "621027"
            }
          ]
        },
        {
          "name": "定西市",
          "code": "621100",
          "sub": [{
            "name": "市辖区",
            "code": "621101"
          },
            {
              "name": "安定区",
              "code": "621102"
            },
            {
              "name": "通渭县",
              "code": "621121"
            },
            {
              "name": "陇西县",
              "code": "621122"
            },
            {
              "name": "渭源县",
              "code": "621123"
            },
            {
              "name": "临洮县",
              "code": "621124"
            },
            {
              "name": "漳县",
              "code": "621125"
            },
            {
              "name": "岷县",
              "code": "621126"
            }
          ]
        },
        {
          "name": "陇南市",
          "code": "621200",
          "sub": [{
            "name": "市辖区",
            "code": "621201"
          },
            {
              "name": "武都区",
              "code": "621202"
            },
            {
              "name": "成县",
              "code": "621221"
            },
            {
              "name": "文县",
              "code": "621222"
            },
            {
              "name": "宕昌县",
              "code": "621223"
            },
            {
              "name": "康县",
              "code": "621224"
            },
            {
              "name": "西和县",
              "code": "621225"
            },
            {
              "name": "礼县",
              "code": "621226"
            },
            {
              "name": "徽县",
              "code": "621227"
            },
            {
              "name": "两当县",
              "code": "621228"
            }
          ]
        },
        {
          "name": "临夏回族自治州",
          "code": "622900",
          "sub": [{
            "name": "临夏市",
            "code": "622901"
          },
            {
              "name": "临夏县",
              "code": "622921"
            },
            {
              "name": "康乐县",
              "code": "622922"
            },
            {
              "name": "永靖县",
              "code": "622923"
            },
            {
              "name": "广河县",
              "code": "622924"
            },
            {
              "name": "和政县",
              "code": "622925"
            },
            {
              "name": "东乡族自治县",
              "code": "622926"
            },
            {
              "name": "积石山保安族东乡族撒拉族自治县",
              "code": "622927"
            }
          ]
        },
        {
          "name": "甘南藏族自治州",
          "code": "623000",
          "sub": [{
            "name": "合作市",
            "code": "623001"
          },
            {
              "name": "临潭县",
              "code": "623021"
            },
            {
              "name": "卓尼县",
              "code": "623022"
            },
            {
              "name": "舟曲县",
              "code": "623023"
            },
            {
              "name": "迭部县",
              "code": "623024"
            },
            {
              "name": "玛曲县",
              "code": "623025"
            },
            {
              "name": "碌曲县",
              "code": "623026"
            },
            {
              "name": "夏河县",
              "code": "623027"
            }
          ]
        }
      ]
    },
    {
      "name": "宁夏",
      "code": "640000",
      "sub": [{
        "name": "银川市",
        "code": "640100",
        "sub": [{
          "name": "市辖区",
          "code": "640101"
        },
          {
            "name": "兴庆区",
            "code": "640104"
          },
          {
            "name": "西夏区",
            "code": "640105"
          },
          {
            "name": "金凤区",
            "code": "640106"
          },
          {
            "name": "永宁县",
            "code": "640121"
          },
          {
            "name": "贺兰县",
            "code": "640122"
          },
          {
            "name": "灵武市",
            "code": "640181"
          }
        ]
      },
        {
          "name": "石嘴山市",
          "code": "640200",
          "sub": [{
            "name": "市辖区",
            "code": "640201"
          },
            {
              "name": "大武口区",
              "code": "640202"
            },
            {
              "name": "惠农区",
              "code": "640205"
            },
            {
              "name": "平罗县",
              "code": "640221"
            }
          ]
        },
        {
          "name": "吴忠市",
          "code": "640300",
          "sub": [{
            "name": "市辖区",
            "code": "640301"
          },
            {
              "name": "利通区",
              "code": "640302"
            },
            {
              "name": "红寺堡区",
              "code": "640303"
            },
            {
              "name": "盐池县",
              "code": "640323"
            },
            {
              "name": "同心县",
              "code": "640324"
            },
            {
              "name": "青铜峡市",
              "code": "640381"
            }
          ]
        },
        {
          "name": "固原市",
          "code": "640400",
          "sub": [{
            "name": "市辖区",
            "code": "640401"
          },
            {
              "name": "原州区",
              "code": "640402"
            },
            {
              "name": "西吉县",
              "code": "640422"
            },
            {
              "name": "隆德县",
              "code": "640423"
            },
            {
              "name": "泾源县",
              "code": "640424"
            },
            {
              "name": "彭阳县",
              "code": "640425"
            }
          ]
        },
        {
          "name": "中卫市",
          "code": "640500",
          "sub": [{
            "name": "市辖区",
            "code": "640501"
          },
            {
              "name": "沙坡头区",
              "code": "640502"
            },
            {
              "name": "中宁县",
              "code": "640521"
            },
            {
              "name": "海原县",
              "code": "640522"
            }
          ]
        }
      ]
    },
    {
      "name": "内蒙古",
      "code": "150000",
      "sub": [{
        "name": "呼和浩特市",
        "code": "150100",
        "sub": [{
          "name": "市辖区",
          "code": "150101"
        },
          {
            "name": "新城区",
            "code": "150102"
          },
          {
            "name": "回民区",
            "code": "150103"
          },
          {
            "name": "玉泉区",
            "code": "150104"
          },
          {
            "name": "赛罕区",
            "code": "150105"
          },
          {
            "name": "土默特左旗",
            "code": "150121"
          },
          {
            "name": "托克托县",
            "code": "150122"
          },
          {
            "name": "和林格尔县",
            "code": "150123"
          },
          {
            "name": "清水河县",
            "code": "150124"
          },
          {
            "name": "武川县",
            "code": "150125"
          }
        ]
      },
        {
          "name": "包头市",
          "code": "150200",
          "sub": [{
            "name": "市辖区",
            "code": "150201"
          },
            {
              "name": "东河区",
              "code": "150202"
            },
            {
              "name": "昆都仑区",
              "code": "150203"
            },
            {
              "name": "青山区",
              "code": "150204"
            },
            {
              "name": "石拐区",
              "code": "150205"
            },
            {
              "name": "白云鄂博矿区",
              "code": "150206"
            },
            {
              "name": "九原区",
              "code": "150207"
            },
            {
              "name": "土默特右旗",
              "code": "150221"
            },
            {
              "name": "固阳县",
              "code": "150222"
            },
            {
              "name": "达尔罕茂明安联合旗",
              "code": "150223"
            }
          ]
        },
        {
          "name": "乌海市",
          "code": "150300",
          "sub": [{
            "name": "市辖区",
            "code": "150301"
          },
            {
              "name": "海勃湾区",
              "code": "150302"
            },
            {
              "name": "海南区",
              "code": "150303"
            },
            {
              "name": "乌达区",
              "code": "150304"
            }
          ]
        },
        {
          "name": "赤峰市",
          "code": "150400",
          "sub": [{
            "name": "市辖区",
            "code": "150401"
          },
            {
              "name": "红山区",
              "code": "150402"
            },
            {
              "name": "元宝山区",
              "code": "150403"
            },
            {
              "name": "松山区",
              "code": "150404"
            },
            {
              "name": "阿鲁科尔沁旗",
              "code": "150421"
            },
            {
              "name": "巴林左旗",
              "code": "150422"
            },
            {
              "name": "巴林右旗",
              "code": "150423"
            },
            {
              "name": "林西县",
              "code": "150424"
            },
            {
              "name": "克什克腾旗",
              "code": "150425"
            },
            {
              "name": "翁牛特旗",
              "code": "150426"
            },
            {
              "name": "喀喇沁旗",
              "code": "150428"
            },
            {
              "name": "宁城县",
              "code": "150429"
            },
            {
              "name": "敖汉旗",
              "code": "150430"
            }
          ]
        },
        {
          "name": "通辽市",
          "code": "150500",
          "sub": [{
            "name": "市辖区",
            "code": "150501"
          },
            {
              "name": "科尔沁区",
              "code": "150502"
            },
            {
              "name": "科尔沁左翼中旗",
              "code": "150521"
            },
            {
              "name": "科尔沁左翼后旗",
              "code": "150522"
            },
            {
              "name": "开鲁县",
              "code": "150523"
            },
            {
              "name": "库伦旗",
              "code": "150524"
            },
            {
              "name": "奈曼旗",
              "code": "150525"
            },
            {
              "name": "扎鲁特旗",
              "code": "150526"
            },
            {
              "name": "霍林郭勒市",
              "code": "150581"
            }
          ]
        },
        {
          "name": "鄂尔多斯市",
          "code": "150600",
          "sub": [{
            "name": "市辖区",
            "code": "150601"
          },
            {
              "name": "东胜区",
              "code": "150602"
            },
            {
              "name": "达拉特旗",
              "code": "150621"
            },
            {
              "name": "准格尔旗",
              "code": "150622"
            },
            {
              "name": "鄂托克前旗",
              "code": "150623"
            },
            {
              "name": "鄂托克旗",
              "code": "150624"
            },
            {
              "name": "杭锦旗",
              "code": "150625"
            },
            {
              "name": "乌审旗",
              "code": "150626"
            },
            {
              "name": "伊金霍洛旗",
              "code": "150627"
            }
          ]
        },
        {
          "name": "呼伦贝尔市",
          "code": "150700",
          "sub": [{
            "name": "市辖区",
            "code": "150701"
          },
            {
              "name": "海拉尔区",
              "code": "150702"
            },
            {
              "name": "扎赉诺尔区",
              "code": "150703"
            },
            {
              "name": "阿荣旗",
              "code": "150721"
            },
            {
              "name": "莫力达瓦达斡尔族自治旗",
              "code": "150722"
            },
            {
              "name": "鄂伦春自治旗",
              "code": "150723"
            },
            {
              "name": "鄂温克族自治旗",
              "code": "150724"
            },
            {
              "name": "陈巴尔虎旗",
              "code": "150725"
            },
            {
              "name": "新巴尔虎左旗",
              "code": "150726"
            },
            {
              "name": "新巴尔虎右旗",
              "code": "150727"
            },
            {
              "name": "满洲里市",
              "code": "150781"
            },
            {
              "name": "牙克石市",
              "code": "150782"
            },
            {
              "name": "扎兰屯市",
              "code": "150783"
            },
            {
              "name": "额尔古纳市",
              "code": "150784"
            },
            {
              "name": "根河市",
              "code": "150785"
            }
          ]
        },
        {
          "name": "巴彦淖尔市",
          "code": "150800",
          "sub": [{
            "name": "市辖区",
            "code": "150801"
          },
            {
              "name": "临河区",
              "code": "150802"
            },
            {
              "name": "五原县",
              "code": "150821"
            },
            {
              "name": "磴口县",
              "code": "150822"
            },
            {
              "name": "乌拉特前旗",
              "code": "150823"
            },
            {
              "name": "乌拉特中旗",
              "code": "150824"
            },
            {
              "name": "乌拉特后旗",
              "code": "150825"
            },
            {
              "name": "杭锦后旗",
              "code": "150826"
            }
          ]
        },
        {
          "name": "乌兰察布市",
          "code": "150900",
          "sub": [{
            "name": "市辖区",
            "code": "150901"
          },
            {
              "name": "集宁区",
              "code": "150902"
            },
            {
              "name": "卓资县",
              "code": "150921"
            },
            {
              "name": "化德县",
              "code": "150922"
            },
            {
              "name": "商都县",
              "code": "150923"
            },
            {
              "name": "兴和县",
              "code": "150924"
            },
            {
              "name": "凉城县",
              "code": "150925"
            },
            {
              "name": "察哈尔右翼前旗",
              "code": "150926"
            },
            {
              "name": "察哈尔右翼中旗",
              "code": "150927"
            },
            {
              "name": "察哈尔右翼后旗",
              "code": "150928"
            },
            {
              "name": "四子王旗",
              "code": "150929"
            },
            {
              "name": "丰镇市",
              "code": "150981"
            }
          ]
        },
        {
          "name": "兴安盟",
          "code": "152200",
          "sub": [{
            "name": "乌兰浩特市",
            "code": "152201"
          },
            {
              "name": "阿尔山市",
              "code": "152202"
            },
            {
              "name": "科尔沁右翼前旗",
              "code": "152221"
            },
            {
              "name": "科尔沁右翼中旗",
              "code": "152222"
            },
            {
              "name": "扎赉特旗",
              "code": "152223"
            },
            {
              "name": "突泉县",
              "code": "152224"
            }
          ]
        },
        {
          "name": "锡林郭勒盟",
          "code": "152500",
          "sub": [{
            "name": "二连浩特市",
            "code": "152501"
          },
            {
              "name": "锡林浩特市",
              "code": "152502"
            },
            {
              "name": "阿巴嘎旗",
              "code": "152522"
            },
            {
              "name": "苏尼特左旗",
              "code": "152523"
            },
            {
              "name": "苏尼特右旗",
              "code": "152524"
            },
            {
              "name": "东乌珠穆沁旗",
              "code": "152525"
            },
            {
              "name": "西乌珠穆沁旗",
              "code": "152526"
            },
            {
              "name": "太仆寺旗",
              "code": "152527"
            },
            {
              "name": "镶黄旗",
              "code": "152528"
            },
            {
              "name": "正镶白旗",
              "code": "152529"
            },
            {
              "name": "正蓝旗",
              "code": "152530"
            },
            {
              "name": "多伦县",
              "code": "152531"
            }
          ]
        },
        {
          "name": "阿拉善盟",
          "code": "152900",
          "sub": [{
            "name": "阿拉善左旗",
            "code": "152921"
          },
            {
              "name": "阿拉善右旗",
              "code": "152922"
            },
            {
              "name": "额济纳旗",
              "code": "152923"
            }
          ]
        }
      ]
    }
  ];
}($);
// jshint ignore: end

/* global $:true */
/* jshint unused:false*/

+ function($) {
  "use strict";

  var defaults;
  var raw = $.rawCitiesData;

  var format = function(data) {
    var result = [];
    for(var i=0;i<data.length;i++) {
      var d = data[i];
      if(/^请选择|市辖区/.test(d.name)) continue;
      result.push(d);
    }
    if(result.length) return result;
    return [];
  };

  var sub = function(data) {
    if(!data.sub) return [{ name: '', code: data.code }];  // 有可能某些县级市没有区
    return format(data.sub);
  };

  var getCities = function(d) {
    for(var i=0;i< raw.length;i++) {
      if(raw[i].code === d || raw[i].name === d) return sub(raw[i]);
    }
    return [];
  };

  var getDistricts = function(p, c) {
    for(var i=0;i< raw.length;i++) {
      if(raw[i].code === p || raw[i].name === p) {
        for(var j=0;j< raw[i].sub.length;j++) {
          if(raw[i].sub[j].code === c || raw[i].sub[j].name === c) {
            return sub(raw[i].sub[j]);
          }
        }
      }
    }
  };

  var parseInitValue = function (val) {
    var p = raw[0], c, d;
    var tokens = val.split(' ');
    raw.map(function (t) {
      if (t.name === tokens[0]) p = t;
    });

    p.sub.map(function (t) {
      if (t.name === tokens[1]) c = t;
    })

    if (tokens[2]) {
      c.sub.map(function (t) {
        if (t.name === tokens[2]) d = t;
      })
    }

    if (d) return [p.code, c.code, d.code];
    return [p.code, c.code];
  }

  $.fn.cityPicker = function(params) {
    params = $.extend({}, defaults, params);
    return this.each(function() {
      var self = this;

      var provincesName = raw.map(function(d) {
        return d.name;
      });
      var provincesCode = raw.map(function(d) {
        return d.code;
      });
      var initCities = sub(raw[0]);
      var initCitiesName = initCities.map(function (c) {
        return c.name;
      });
      var initCitiesCode = initCities.map(function (c) {
        return c.code;
      });
      var initDistricts = sub(raw[0].sub[0]);

      var initDistrictsName = initDistricts.map(function (c) {
        return c.name;
      });
      var initDistrictsCode = initDistricts.map(function (c) {
        return c.code;
      });

      var currentProvince = provincesName[0];
      var currentCity = initCitiesName[0];
      var currentDistrict = initDistrictsName[0];

      var cols = [
        {
          displayValues: provincesName,
          values: provincesCode,
          cssClass: "col-province"
        },
        {
          displayValues: initCitiesName,
          values: initCitiesCode,
          cssClass: "col-city"
        }
      ];

      if(params.showDistrict) cols.push({
        values: initDistrictsCode,
        displayValues: initDistrictsName,
        cssClass: "col-district"
      });

      var config = {

        cssClass: "city-picker",
        rotateEffect: false,  //为了性能
        formatValue: function (p, values, displayValues) {
          return displayValues.join(' ');
        },
        onChange: function (picker, values, displayValues) {
          var newProvince = picker.cols[0].displayValue;
          var newCity;
          if(newProvince !== currentProvince) {
            var newCities = getCities(newProvince);
            newCity = newCities[0].name;
            var newDistricts = getDistricts(newProvince, newCity);
            picker.cols[1].replaceValues(newCities.map(function (c) {
              return c.code;
            }), newCities.map(function (c) {
              return c.name;
            }));
            if(params.showDistrict) picker.cols[2].replaceValues(newDistricts.map(function (d) {
              return d.code;
            }), newDistricts.map(function (d) {
              return d.name;
            }));
            currentProvince = newProvince;
            currentCity = newCity;
            picker.updateValue();
            return false; // 因为数据未更新完，所以这里不进行后序的值的处理
          } else {
            if(params.showDistrict) {
              newCity = picker.cols[1].displayValue;
              if(newCity !== currentCity) {
                var districts = getDistricts(newProvince, newCity);
                picker.cols[2].replaceValues(districts.map(function (d) {
                  return d.code;
                }), districts.map(function (d) {
                  return d.name;
                }));
                currentCity = newCity;
                picker.updateValue();
                return false; // 因为数据未更新完，所以这里不进行后序的值的处理
              }
            }
          }
          //如果最后一列是空的，那么取倒数第二列
          var len = (values[values.length-1] ? values.length - 1 : values.length - 2)
          $(self).attr('data-code', values[len]);
          $(self).attr('data-codes', values.join(','));
          if (params.onChange) {
            params.onChange.call(self, picker, values, displayValues);
          }
        },

        cols: cols
      };

      if(!this) return;
      var p = $.extend({}, params, config);
      //计算value
      var val = $(this).val();
      if (!val) val = '北京 北京市 东城区';
      currentProvince = val.split(" ")[0];
      currentCity = val.split(" ")[1];
      currentDistrict= val.split(" ")[2];
      if(val) {
        p.value = parseInitValue(val);
        if(p.value[0]) {
          var cities = getCities(p.value[0]);
          p.cols[1].values = cities.map(function (c) {
            return c.code;
          });
          p.cols[1].displayValues = cities.map(function (c) {
            return c.name;
          });
        }

        if(p.value[1]) {
          if (params.showDistrict) {
            var dis = getDistricts(p.value[0], p.value[1]);
            p.cols[2].values = dis.map(function (d) {
              return d.code;
            });
            p.cols[2].displayValues = dis.map(function (d) {
              return d.name;
            });
          }
        } else {
          if (params.showDistrict) {
            var dis = getDistricts(p.value[0], p.cols[1].values[0]);
            p.cols[2].values = dis.map(function (d) {
              return d.code;
            });
            p.cols[2].displayValues = dis.map(function (d) {
              return d.name;
            });
          }
        }
      }
      $(this).picker(p);
    });
  };

  defaults = $.fn.cityPicker.prototype.defaults = {
    showDistrict: true //是否显示地区选择
  };

}($);