module Logic(
  input [3:0] a,
  input [3:0] b,
  output [3:0] result
);

  wire [3:0] temp;

  assign temp = a & b;
  assign result = temp;

endmodule
