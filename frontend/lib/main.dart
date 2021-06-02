import 'package:flutter/material.dart';
import 'package:frontend/mow_plow_home.dart';

void main() => runApp(MowPlow());

class MowPlow extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MowPlow',
      theme: ThemeData(
        primarySwatch: Colors.green,
      ),
      debugShowCheckedModeBanner: false,
      home: new MowPlowHome(),
      );
  }
}