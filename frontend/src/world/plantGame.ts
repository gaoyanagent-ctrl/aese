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
  | "project_baselined";

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
};

export type PlantGameStageDefinition = {
  key: PlantGameStage;
  label: string;
  location: "headquarters" | "city" | "site" | "boardroom";
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
    dialogue: "项目和 WBS 已进入 IAOS 权威账。接下来的合同、施工和工程财务都必须引用这条基线。",
    actionLabel: "查看项目基线档案",
    anchor: "plant-task-project",
  },
];

export function derivePlantGameStage(facts: PlantGameFacts): PlantGameStageDefinition {
  if (facts.hasActiveProject) return PLANT_GAME_STAGES[9];
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
