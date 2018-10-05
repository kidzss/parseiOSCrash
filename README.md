# parseiOSCrash

使用symbolicatecrash 批量解析 iOS Crash，使用golang编写，可以把对应的dSYM文件，*.app文件，批量的  *.crash文件，或者  *.ips文件解析成对应的名字的 *.log文件

1.可以自动把symbolicatecrash copy到当前的目录

2.把文件准备好批量执行

3.代码可以切换 输入到文件，可以打印到控制台

```
symbolicatecrash path：/Applications/Xcode.app/Contents/SharedFrameworks/DVTFoundation.framework/Versions/A/Resources/symbolicatecrash
```