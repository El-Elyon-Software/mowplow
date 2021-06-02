
import 'package:flutter/material.dart';

class MowPlowHome extends StatefulWidget {
  @override
  _MowPlowHome createState() => _MowPlowHome();
}

class _MowPlowHome extends State<MowPlowHome> {
  var _unscheduledCount = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text("Mow Plow"),
        ),
        body: new Row(
          children: [
            ElevatedButton(
              style: ElevatedButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(5),
                ),
                primary: Colors.green,
              ),
              child: Container(
                width: 75,
                height: 100,
                alignment: Alignment.center,
                
                child: Text(
                  _unscheduledCount.toString(),
                  style: TextStyle(fontSize: 45),
                ),
              ),
              onPressed: () {},
            ),
            ElevatedButton(
              style: ElevatedButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(5),
                ),
                primary: Colors.green,
                image: 
              ),
              child: Container(
                width: 75,
                height: 100,
                alignment: Alignment.center,
                
                child: Text(
                  'Current Schedule',
                  style: TextStyle(fontSize: 15),
                ),
              ),
              onPressed: () {},
            )
          ],
        ));
  }
}
