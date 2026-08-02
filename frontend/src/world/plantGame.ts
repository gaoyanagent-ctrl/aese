export type PlantGameStage =
  | "requirement"
  | "proposal"
  | "investigation"
  | "comparison"
  | "governance"
  | "selected"
  | "site_control"
  | "project_planning"
  | "project_approval"
  | "project_baselined"
  | "contract_sourcing"
  | "contract_world"
  | "contract_comparison"
  | "contract_approval"
  | "contract_awarded"
  | "construction_start"
  | "construction_world"
  | "milestone_acceptance"
  | "milestone_accepted";

export type PlantGameFacts = {
  hasRequirement: boolean;
  proposalCount: number;
  adoptedReviewCount: number;
  investigationCount: number;
  observedCount: number;
  hasRecommendation: boolean;
  hasDecision: boolean;
  hasSiteControlRequest: boolean;
  hasSiteControl: boolean;
  hasProjectPlan: boolean;
  projectApprovalStatus: string;
  hasActiveProject: boolean;
  hasContractRFQ: boolean;
  hasContractBids: boolean;
  hasContractRecommendation: boolean;
  contractApprovalStatus: string;
  hasAwardedContract: boolean;
  hasConstructionExecution: boolean;
  hasConstructionObservation: boolean;
  hasMilestoneAcceptance: boolean;
};

export type PlantGameStageDefinition = {
  key: PlantGameStage;
  label: string;
  location: "headquarters" | "city" | "site" | "boardroom" | "contractor" | "construction";
  npc: string;
  npcRole: string;
  mission: string;
  dialogue: string;
  actionLabel: string;
  anchor: string;
};

export const PLANT_GAME_STAGES: PlantGameStageDefinition[] = [
  {
    key: "requirement",
    label: "定义设施需求",
    location: "headquarters",
    npc: "纪元",
    npcRole: "设施规划 Agent",
    mission: "在总部规划室明确工厂需要什么",
    dialogue: "我会读取本企业已过账现金和已批准预算。请你决定区域、面积、电力、时间和可承受投资边界。",
    actionLabel: "与规划 Agent 制定需求",
    anchor: "plant-task-requirement",
  },
  {
    key: "proposal",
    label: "审阅候选方案",
    location: "city",
    npc: "纪元",
    npcRole: "设施规划 Agent",
    mission: "在产业地图上比较候选场址",
    dialogue: "候选只是有依据的建议，不是报价或批准。请选择值得进入真实调研的方案，并说明理由。",
    actionLabel: "查看候选与作出评审",
    anchor: "plant-task-proposals",
  },
  {
    key: "investigation",
    label: "等待外部调研",
    location: "site",
    npc: "园区顾问",
    npcRole: "World 外部参与者",
    mission: "前往候选园区核验现实条件",
    dialogue: "权属、正式报价、可用面积、电力和许可必须由外部事实证明。IAOS 不会把 Agent 估算当成现场结果。",
    actionLabel: "进入园区调研工作项",
    anchor: "plant-task-investigation",
  },
  {
    key: "comparison",
    label: "比较外部事实",
    location: "headquarters",
    npc: "周衡",
    npcRole: "项目与财务负责人",
    mission: "回到总部形成可解释推荐",
    dialogue: "先检查投资、日期、面积、电力、权属和许可硬约束，再调整经营权重。高分不能覆盖硬约束失败。",
    actionLabel: "比较事实并提交推荐",
    anchor: "plant-task-comparison",
  },
  {
    key: "governance",
    label: "治理审批",
    location: "boardroom",
    npc: "林岚",
    npcRole: "治理与审批协调人",
    mission: "在治理会议室等待有权主体决定",
    dialogue: "推荐已经冻结为审批事项。审批流决定由谁审阅；Agent、提交人和游戏画面都不能自行批准。",
    actionLabel: "查看审批与正式选址",
    anchor: "plant-task-governance",
  },
  {
    key: "selected",
    label: "场址已确定",
    location: "site",
    npc: "顾远",
    npcRole: "工厂项目负责人",
    mission: "确认场址并准备建设项目",
    dialogue: "正式场址已经由受治理能力落地。下一步将建立项目、WBS、合同、施工、付款和验收闭环。",
    actionLabel: "查看正式选址证据",
    anchor: "plant-task-governance",
  },
  {
    key: "site_control",
    label: "取得场地控制",
    location: "site",
    npc: "园区权利方",
    npcRole: "World 外部参与者",
    mission: "让正式场址从治理决定变成真实可占有场地",
    dialogue: "选址批准不是交付。请发起场地控制请求；我会基于已签协议、交付记录和占有权限返回可信 Observation。",
    actionLabel: "办理协议与场地交付",
    anchor: "plant-task-site-control",
  },
  {
    key: "project_planning",
    label: "建立项目基线",
    location: "headquarters",
    npc: "纪元",
    npcRole: "设施项目 Agent",
    mission: "在项目办公室准备设施项目与 WBS",
    dialogue: "场地已经交付。我会准备交付策略、预算、日期和专业 WBS；你只需要选择方案并确认管理边界。",
    actionLabel: "让 Agent 准备项目方案",
    anchor: "plant-task-project",
  },
  {
    key: "project_approval",
    label: "审批项目基线",
    location: "boardroom",
    npc: "林岚",
    npcRole: "治理与审批协调人",
    mission: "审阅设施项目与 WBS 基线",
    dialogue: "项目草案已经冻结为审批事项。审批流决定有权主体；批准后项目负责人还要显式激活基线。",
    actionLabel: "进入项目基线审批",
    anchor: "plant-task-project",
  },
  {
    key: "project_baselined",
    label: "项目基线已激活",
    location: "headquarters",
    npc: "顾远",
    npcRole: "工厂项目负责人",
    mission: "查验项目与 WBS 档案",
    dialogue: "项目和 WBS 已进入 IAOS 权威账。现在选择一个专业工作包，启动受治理的承包商寻源。",
    actionLabel: "进入工程采购寻源",
    anchor: "plant-task-contract",
  },
  {
    key: "contract_sourcing",
    label: "发布工程采购邀请",
    location: "headquarters",
    npc: "顾远",
    npcRole: "工厂项目负责人",
    mission: "从已批准 WBS 选择采购包",
    dialogue: "采购包、合同上限和交付日期全部来自已激活项目基线。你只需选择包和寻源策略。",
    actionLabel: "选择采购包并发布 RFQ",
    anchor: "plant-task-contract",
  },
  {
    key: "contract_world",
    label: "收取可信投标",
    location: "contractor",
    npc: "市场协调人",
    npcRole: "World 承包商市场",
    mission: "让外部承包商提交密封报价",
    dialogue: "承包商、报价、资质、工期和证据由 World 生成。你只需要确认接收本轮正式投标。",
    actionLabel: "收取正式投标",
    anchor: "plant-task-contract",
  },
  {
    key: "contract_comparison",
    label: "Agent 比选推荐",
    location: "headquarters",
    npc: "纪元",
    npcRole: "工程采购评审 Agent",
    mission: "比较可信报价并形成授予建议",
    dialogue: "我只使用 World 已提交的正式投标，解释成本、工期、质保与履约权衡；最终仍由你确认并提交审批。",
    actionLabel: "让 Agent 评审投标",
    anchor: "plant-task-contract",
  },
  {
    key: "contract_approval",
    label: "审批合同授予",
    location: "boardroom",
    npc: "林岚",
    npcRole: "治理与审批协调人",
    mission: "审阅合同承诺与中选依据",
    dialogue: "审批流决定有权主体。批准后仍需项目负责人显式归档合同，系统不会自动付款或生成发票。",
    actionLabel: "进入合同授予审批",
    anchor: "plant-task-contract",
  },
  {
    key: "contract_awarded",
    label: "正式合同已归档",
    location: "headquarters",
    npc: "顾远",
    npcRole: "工厂项目负责人",
    mission: "查验工程合同与承诺金额",
    dialogue: "正式合同已进入 IAOS 权威账。后续施工进度、变更、发票和付款只能引用这份合同。",
    actionLabel: "启动合同施工包",
    anchor: "plant-task-construction",
  },
  {
    key: "construction_start", label: "启动施工包", location: "headquarters", npc: "顾远", npcRole: "工厂项目负责人",
    mission: "按正式合同启动施工", dialogue: "合同只代表承诺。点击确认后，我会向施工现场发出受治理开工意图；不会预填工程进度。", actionLabel: "确认启动施工", anchor: "plant-task-construction",
  },
  {
    key: "construction_world", label: "现场施工与检查", location: "construction", npc: "现场经理", npcRole: "World 施工承包商",
    mission: "推进施工并提交可信现场报告", dialogue: "进度、质量、安全和现场证据由施工 World 生成。你只需推进剧情，不填写工程数据。", actionLabel: "推进施工并提交检查报告", anchor: "plant-task-construction",
  },
  {
    key: "milestone_acceptance", label: "验收施工里程碑", location: "construction", npc: "顾远", npcRole: "工厂项目负责人",
    mission: "核验现场事实并明确验收", dialogue: "现场报告已送达。只有进度、质量、安全和遗留项门全部通过，才能确认里程碑；验收不会自动付款。", actionLabel: "核验并验收里程碑", anchor: "plant-task-construction",
  },
  {
    key: "milestone_accepted", label: "里程碑已验收", location: "headquarters", npc: "顾远", npcRole: "工厂项目负责人",
    mission: "查验施工与验收档案", dialogue: "施工事实和验收已经分别入账，付款状态仍为未申请。下一步将建立工程发票、应付和付款治理。", actionLabel: "查看施工验收档案", anchor: "plant-task-construction",
  },
];

export function derivePlantGameStage(facts: PlantGameFacts): PlantGameStageDefinition {
  if (facts.hasMilestoneAcceptance) return PLANT_GAME_STAGES[18];
  if (facts.hasConstructionObservation) return PLANT_GAME_STAGES[17];
  if (facts.hasConstructionExecution) return PLANT_GAME_STAGES[16];
  if (facts.hasAwardedContract) return PLANT_GAME_STAGES[15];
  if (facts.hasContractRecommendation) return PLANT_GAME_STAGES[13];
  if (facts.hasContractBids) return PLANT_GAME_STAGES[12];
  if (facts.hasContractRFQ) return PLANT_GAME_STAGES[11];
  if (facts.hasActiveProject) return PLANT_GAME_STAGES[10];
  if (facts.hasProjectPlan && facts.projectApprovalStatus !== "rejected") return PLANT_GAME_STAGES[8];
  if (facts.projectApprovalStatus === "rejected") return PLANT_GAME_STAGES[7];
  if (facts.hasSiteControl) return PLANT_GAME_STAGES[7];
  // A formal decision completes the read-only "selected" checkpoint and must
  // immediately expose the next human/World task. Waiting for a control request
  // before showing the request UI would make that first command unreachable.
  if (facts.hasDecision) return PLANT_GAME_STAGES[6];
  if (facts.hasRecommendation) return PLANT_GAME_STAGES[4];
  if (facts.observedCount > 0) return PLANT_GAME_STAGES[3];
  if (facts.investigationCount > 0 || facts.adoptedReviewCount > 0) return PLANT_GAME_STAGES[2];
  if (facts.proposalCount > 0 || facts.hasRequirement) return PLANT_GAME_STAGES[1];
  return PLANT_GAME_STAGES[0];
}
