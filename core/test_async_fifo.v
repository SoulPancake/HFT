module AsyncFIFOTest(
  input [0:0] wr_clk,
  input [0:0] wr_rst,
  input [0:0] rd_clk,
  input [0:0] rd_rst,
  input [15:0] test_fifo_wr_data,
  input [0:0] test_fifo_wr_en,
  input [0:0] test_fifo_rd_en,
  output [0:0] test_fifo_wr_full,
  output [15:0] test_fifo_rd_data,
  output [0:0] test_fifo_rd_empty,
  output [15:0] wr_data_out,
  output [15:0] rd_data_out,
  output [0:0] wr_full_out,
  output [0:0] rd_empty_out
);

  reg [5:0] test_fifo_wr_ptr;
  reg [5:0] test_fifo_rd_ptr;
  reg [5:0] test_fifo_wr_ptr_sync_sync_stage0;
  reg [5:0] test_fifo_wr_ptr_sync_sync_stage1;
  reg [5:0] test_fifo_rd_ptr_sync_sync_stage0;
  reg [5:0] test_fifo_rd_ptr_sync_sync_stage1;
  reg [15:0] test_fifo_mem [0:31];

  assign wr_data_out = test_fifo_wr_data;
  assign rd_data_out = test_fifo_rd_data;
  assign wr_full_out = test_fifo_wr_full;
  assign rd_empty_out = test_fifo_rd_empty;

  // Clock Domains:
  // - write: clk=wr_clk, rst=wr_rst
  // - read: clk=rd_clk, rst=rd_rst

  always @(posedge rd_clk) test_fifo_wr_ptr_sync_sync_stage0 <= test_fifo_wr_ptr;

  always @(posedge rd_clk) test_fifo_wr_ptr_sync_sync_stage1 <= test_fifo_wr_ptr_sync_sync_stage0;

  always @(posedge wr_clk) test_fifo_rd_ptr_sync_sync_stage0 <= test_fifo_rd_ptr;

  always @(posedge wr_clk) test_fifo_rd_ptr_sync_sync_stage1 <= test_fifo_rd_ptr_sync_sync_stage0;

  // AsyncFIFO test_fifo implementation

  // Memory: test_fifo_mem

  // Write enable: test_fifo_wr_en

  // Read enable: test_fifo_rd_en

  // Write pointer sync: test_fifo_wr_ptr_sync_sync_stage1

  // Read pointer sync: test_fifo_rd_ptr_sync_sync_stage1

endmodule
