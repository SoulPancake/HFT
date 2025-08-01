module TopModule(
  input [0:0] clk,
  output [15:0] output_data
);

  assign output_data = sub_inst_data_out;

  SubModule #(
    .WIDTH(16)
  ) sub_inst (
    .clk(clk),
    .data_out(sub_inst_data_out)
  );

endmodule
