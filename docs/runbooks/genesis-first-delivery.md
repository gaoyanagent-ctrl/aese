# M13 Genesis First Delivery Runbook

访问 `/#world-industrialization`，进入“首次商业交付 Campaign”，单步 12 次。终态必须显示 `Genesis cycle closed`、三批 9000/2700/300、累计接受 12,000、AR 0、实际成本与项目毛利。

运行全仓 Go test/vet、pack validate、前端 test/typecheck/build/Playwright。IAOS `/genesis/delivery/actions` 的 accept/collect/cost/close 必须引用 World evidence；发运、客户接受、发票、应收、银行到账和核销严格分离。
