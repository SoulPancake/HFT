module ClockDomainTest(
  input [0:0] clk1,
  input [0:0] rst1,
  input [0:0] clk2,
  input [0:0] rst2
);

  wire [7:0] sig1;
  wire [7:0] sig2;
  reg [7:0] cross_sync_sync_stage0;
  reg [7:0] cross_sync_sync_stage1;

  // Clock Domains:
  // - fast_domain: clk=clk1, rst=rst1, freq=200000000 Hz
  // - slow_domain: clk=clk2, rst=rst2, freq=50000000 Hz

  always @(posedge clk2) cross_sync_sync_stage0 <= sig1;

  always @(posedge clk2) cross_sync_sync_stage1 <= cross_sync_sync_stage0;

endmodule
