module ComplexSystem(
  input [0:0] cpu_clk,
  input [0:0] cpu_rst,
  input [0:0] ddr_clk,
  input [0:0] ddr_rst,
  input [0:0] bus_req_0,
  input [0:0] bus_req_1,
  input [0:0] bus_req_2,
  input [0:0] bus_req_3,
  input [0:0] memory_req_0,
  input [0:0] memory_req_1,
  input [31:0] main_fifo_wr_data,
  input [0:0] main_fifo_wr_en,
  input [0:0] main_fifo_rd_en,
  output [0:0] bus_grant_0,
  output [0:0] bus_grant_1,
  output [0:0] bus_grant_2,
  output [0:0] bus_grant_3,
  output [0:0] memory_grant_0,
  output [0:0] memory_grant_1,
  output [0:0] main_fifo_wr_full,
  output [31:0] main_fifo_rd_data,
  output [0:0] main_fifo_rd_empty,
  output [31:0] data_out,
  output [7:0] status_out
);

  wire [31:0] cpu_data;
  reg [0:0] memory_counter;
  reg [31:0] cpu_to_ddr_sync_stage0;
  reg [31:0] cpu_to_ddr_sync_stage1;
  reg [31:0] cpu_to_ddr_sync_stage2;
  reg [7:0] main_fifo_wr_ptr;
  reg [7:0] main_fifo_rd_ptr;
  reg [7:0] main_fifo_wr_ptr_sync_sync_stage0;
  reg [7:0] main_fifo_wr_ptr_sync_sync_stage1;
  reg [7:0] main_fifo_rd_ptr_sync_sync_stage0;
  reg [7:0] main_fifo_rd_ptr_sync_sync_stage1;
  reg [31:0] main_fifo_mem [0:127];

  assign bus_grant_3 = bus_req_3;
  assign bus_grant_2 = bus_req_2 && !(bus_req_3);
  assign bus_grant_1 = bus_req_1 && !(bus_req_2 || bus_req_3);
  assign bus_grant_0 = bus_req_0 && !(bus_req_1 || bus_req_2 || bus_req_3);
  assign data_out = cpu_to_ddr_sync_stage2;
  assign status_out = {bus_grant_0, memory_grant_0, main_fifo_wr_full, main_fifo_rd_empty, {4{0}}};

  // Clock Domains:
  // - cpu: clk=cpu_clk, rst=cpu_rst, freq=100000000 Hz
  // - ddr: clk=ddr_clk, rst=ddr_rst, freq=200000000 Hz

  // Mutex: bus (priority arbitration)
  // Request[0]: bus_req_0 -> Grant[0]: bus_grant_0
  // Request[1]: bus_req_1 -> Grant[1]: bus_grant_1
  // Request[2]: bus_req_2 -> Grant[2]: bus_grant_2
  // Request[3]: bus_req_3 -> Grant[3]: bus_grant_3

  // Mutex: memory (round_robin arbitration)
  // Request[0]: memory_req_0 -> Grant[0]: memory_grant_0
  // Request[1]: memory_req_1 -> Grant[1]: memory_grant_1

  always @(posedge cpu_clk) begin

    memory_grant_0 <= (memory_counter == 0) && memory_req_0;

    memory_grant_1 <= (memory_counter == 1) && memory_req_1;

    if (memory_grant_0 || memory_grant_1) memory_counter <= (memory_counter + 1) % 2;

  end

  always @(posedge ddr_clk) cpu_to_ddr_sync_stage0 <= cpu_data;

  always @(posedge ddr_clk) cpu_to_ddr_sync_stage1 <= cpu_to_ddr_sync_stage0;

  always @(posedge ddr_clk) cpu_to_ddr_sync_stage2 <= cpu_to_ddr_sync_stage1;

  always @(posedge ddr_clk) main_fifo_wr_ptr_sync_sync_stage0 <= main_fifo_wr_ptr;

  always @(posedge ddr_clk) main_fifo_wr_ptr_sync_sync_stage1 <= main_fifo_wr_ptr_sync_sync_stage0;

  always @(posedge cpu_clk) main_fifo_rd_ptr_sync_sync_stage0 <= main_fifo_rd_ptr;

  always @(posedge cpu_clk) main_fifo_rd_ptr_sync_sync_stage1 <= main_fifo_rd_ptr_sync_sync_stage0;

  // AsyncFIFO main_fifo implementation

  // Memory: main_fifo_mem

  // Write enable: main_fifo_wr_en

  // Read enable: main_fifo_rd_en

  // Write pointer sync: main_fifo_wr_ptr_sync_sync_stage1

  // Read pointer sync: main_fifo_rd_ptr_sync_sync_stage1

endmodule
