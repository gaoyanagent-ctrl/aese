import preview from "../../../scenario-packs/hctm/stories/order-expedite-01/preview.json";
import { StaticScenarioDataSource } from "./dataSource";

export const DEFAULT_SCENARIO_KEY = "order-expedite-01";

export const previewScenarioDataSource = new StaticScenarioDataSource({
  [DEFAULT_SCENARIO_KEY]: preview,
});
