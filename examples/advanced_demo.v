module AdvancedHFTDemo(
  input [0:0] cpu_clk,
  input [0:0] cpu_rst,
  input [0:0] mem_clk,
  input [0:0] mem_rst,
  input [0:0] bus_arbiter_req_0,
  input [0:0] bus_arbiter_req_1,
  input [0:0] bus_arbiter_req_2,
  input [0:0] memory_arbiter_req_0,
  input [0:0] memory_arbiter_req_1,
  input [31:0] cpu_mem_fifo_wr_data,
  input [0:0] cpu_mem_fifo_wr_en,
  input [0:0] cpu_mem_fifo_rd_en,
  input [31:0] data_in,
  input [4:0] addr_in,
  input [0:0] enable_in,
  output [0:0] bus_arbiter_grant_0,
  output [0:0] bus_arbiter_grant_1,
  output [0:0] bus_arbiter_grant_2,
  output [0:0] memory_arbiter_grant_0,
  output [0:0] memory_arbiter_grant_1,
  output [0:0] cpu_mem_fifo_wr_full,
  output [31:0] cpu_mem_fifo_rd_data,
  output [0:0] cpu_mem_fifo_rd_empty,
  output [31:0] data_out,
  output [7:0] status_out
);

  wire [31:0] cpu_data;
  reg [0:0] memory_arbiter_counter;
  reg [31:0] cpu_to_mem_sync_stage0;
  reg [31:0] cpu_to_mem_sync_stage1;
  reg [6:0] cpu_mem_fifo_wr_ptr;
  reg [6:0] cpu_mem_fifo_rd_ptr;
  reg [6:0] cpu_mem_fifo_wr_ptr_sync_sync_stage0;
  reg [6:0] cpu_mem_fifo_wr_ptr_sync_sync_stage1;
  reg [6:0] cpu_mem_fifo_rd_ptr_sync_sync_stage0;
  reg [6:0] cpu_mem_fifo_rd_ptr_sync_sync_stage1;
  reg [31:0] cpu_mem_fifo_mem [0:63];

  assign bus_arbiter_grant_2 = bus_arbiter_req_2;
  assign bus_arbiter_grant_1 = bus_arbiter_req_1 && !(bus_arbiter_req_2);
  assign bus_arbiter_grant_0 = bus_arbiter_req_0 && !(bus_arbiter_req_1 || bus_arbiter_req_2);
  assign cpu_data = data_in;
  assign data_out = cpu_to_mem_sync_stage1;
  assign status_out = {bus_arbiter_grant_0, memory_arbiter_grant_0, cpu_mem_fifo_wr_full, cpu_mem_fifo_rd_empty, enable_in, {3{0}}};

  // Clock Domains:
  // - cpu_domain: clk=cpu_clk, rst=cpu_rst, freq=100000000 Hz
  // - mem_domain: clk=mem_clk, rst=mem_rst, freq=200000000 Hz

  // Mutex: bus_arbiter (priority arbitration)
  // Request[0]: bus_arbiter_req_0 -> Grant[0]: bus_arbiter_grant_0
  // Request[1]: bus_arbiter_req_1 -> Grant[1]: bus_arbiter_grant_1
  // Request[2]: bus_arbiter_req_2 -> Grant[2]: bus_arbiter_grant_2

  // Mutex: memory_arbiter (round_robin arbitration)
  // Request[0]: memory_arbiter_req_0 -> Grant[0]: memory_arbiter_grant_0
  // Request[1]: memory_arbiter_req_1 -> Grant[1]: memory_arbiter_grant_1

  always @(posedge cpu_clk) begin

    memory_arbiter_grant_0 <= (memory_arbiter_counter == 0) && memory_arbiter_req_0;

    memory_arbiter_grant_1 <= (memory_arbiter_counter == 1) && memory_arbiter_req_1;

    if (memory_arbiter_grant_0 || memory_arbiter_grant_1) memory_arbiter_counter <= (memory_arbiter_counter + 1) % 2;

  end

  always @(posedge mem_clk) cpu_to_mem_sync_stage0 <= cpu_data;

  always @(posedge mem_clk) cpu_to_mem_sync_stage1 <= cpu_to_mem_sync_stage0;

  always @(posedge mem_clk) cpu_mem_fifo_wr_ptr_sync_sync_stage0 <= cpu_mem_fifo_wr_ptr;

  always @(posedge mem_clk) cpu_mem_fifo_wr_ptr_sync_sync_stage1 <= cpu_mem_fifo_wr_ptr_sync_sync_stage0;

  always @(posedge cpu_clk) cpu_mem_fifo_rd_ptr_sync_sync_stage0 <= cpu_mem_fifo_rd_ptr;

  always @(posedge cpu_clk) cpu_mem_fifo_rd_ptr_sync_sync_stage1 <= cpu_mem_fifo_rd_ptr_sync_sync_stage0;

  // AsyncFIFO cpu_mem_fifo implementation

  // Memory: cpu_mem_fifo_mem

  // Write enable: cpu_mem_fifo_wr_en

  // Read enable: cpu_mem_fifo_rd_en

  // Write pointer sync: cpu_mem_fifo_wr_ptr_sync_sync_stage1

  // Read pointer sync: cpu_mem_fifo_rd_ptr_sync_sync_stage1

endmodule
