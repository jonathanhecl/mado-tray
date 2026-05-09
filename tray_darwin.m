#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

extern void madoTrayShow(void);
extern void madoTrayExit(void);

@interface MadoTrayMenuTarget : NSObject
- (void)showWindow:(id)sender;
- (void)exitApp:(id)sender;
@end

@implementation MadoTrayMenuTarget
- (void)showWindow:(id)sender {
  [NSApp activateIgnoringOtherApps:YES];
  madoTrayShow();
}

- (void)exitApp:(id)sender {
  madoTrayExit();
}
@end

static NSStatusItem *madoTrayStatusItem;
static MadoTrayMenuTarget *madoTrayTarget;
static NSMenuItem *madoTrayShowItem;
static NSMenuItem *madoTrayExitItem;
static NSString *madoTrayLocale = @"en";

static NSString *MadoTrayShowLabel(void) {
  if ([madoTrayLocale isEqualToString:@"es"]) {
    return @"Mostrar";
  }
  return @"Show";
}

static NSString *MadoTrayExitLabel(void) {
  if ([madoTrayLocale isEqualToString:@"es"]) {
    return @"Salir de Mado-Tray";
  }
  return @"Exit Mado-Tray";
}

static NSImage *MadoTrayIcon(void) {
  const CGFloat size = 18.0;
  NSImage *image = [[NSImage alloc] initWithSize:NSMakeSize(size, size)];

  [image lockFocus];
  [[NSColor blackColor] setStroke];
  [[NSColor blackColor] setFill];

  NSBezierPath *window = [NSBezierPath bezierPathWithRoundedRect:NSMakeRect(2.0, 3.0, 14.0, 12.0)
                                                        xRadius:3.0
                                                        yRadius:3.0];
  [window setLineWidth:1.7];
  [window stroke];

  NSBezierPath *header = [NSBezierPath bezierPath];
  [header moveToPoint:NSMakePoint(3.0, 11.5)];
  [header lineToPoint:NSMakePoint(15.0, 11.5)];
  [header setLineWidth:1.2];
  [header stroke];

  [[NSBezierPath bezierPathWithOvalInRect:NSMakeRect(4.4, 12.7, 1.5, 1.5)] fill];
  [[NSBezierPath bezierPathWithOvalInRect:NSMakeRect(6.8, 12.7, 1.5, 1.5)] fill];

  NSBezierPath *row1 = [NSBezierPath bezierPathWithRoundedRect:NSMakeRect(6.0, 8.5, 7.8, 1.6)
                                                       xRadius:0.8
                                                       yRadius:0.8];
  [row1 fill];

  NSBezierPath *row2 = [NSBezierPath bezierPathWithRoundedRect:NSMakeRect(6.0, 5.8, 6.0, 1.6)
                                                       xRadius:0.8
                                                       yRadius:0.8];
  [row2 fill];

  [[NSBezierPath bezierPathWithOvalInRect:NSMakeRect(4.0, 8.4, 1.8, 1.8)] fill];
  [[NSBezierPath bezierPathWithOvalInRect:NSMakeRect(4.0, 5.7, 1.8, 1.8)] fill];

  [image unlockFocus];
  image.template = YES;
  return image;
}

static void MadoTrayEnsureTarget(void) {
  if (madoTrayTarget == nil) {
    madoTrayTarget = [[MadoTrayMenuTarget alloc] init];
  }
}

static void MadoTrayUpdateMenuTexts(void) {
  if (madoTrayShowItem == nil || madoTrayExitItem == nil) {
    return;
  }

  madoTrayShowItem.title = MadoTrayShowLabel();
  madoTrayExitItem.title = MadoTrayExitLabel();
}

void MadoTrayCreate(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    MadoTrayEnsureTarget();
  });
}

void MadoTrayShow(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    MadoTrayEnsureTarget();

    if (madoTrayStatusItem != nil) {
      return;
    }

    madoTrayStatusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    [madoTrayStatusItem retain];
    madoTrayStatusItem.autosaveName = @"com.jonathanhecl.mado-tray.status-item";
    if (@available(macOS 10.12, *)) {
      madoTrayStatusItem.behavior = 0;
    }
    madoTrayStatusItem.button.toolTip = @"Mado-Tray";

    NSImage *image = MadoTrayIcon();
    if (image != nil) {
      madoTrayStatusItem.button.image = image;
    } else {
      madoTrayStatusItem.button.title = @"Mado";
    }
    madoTrayStatusItem.button.enabled = YES;
    madoTrayStatusItem.button.hidden = NO;

    NSMenu *menu = [[NSMenu alloc] initWithTitle:@"Mado-Tray"];
    madoTrayShowItem = [[NSMenuItem alloc] initWithTitle:MadoTrayShowLabel()
                                                   action:@selector(showWindow:)
                                            keyEquivalent:@""];
    madoTrayShowItem.target = madoTrayTarget;
    [menu addItem:madoTrayShowItem];

    [menu addItem:[NSMenuItem separatorItem]];

    madoTrayExitItem = [[NSMenuItem alloc] initWithTitle:MadoTrayExitLabel()
                                                   action:@selector(exitApp:)
                                            keyEquivalent:@""];
    madoTrayExitItem.target = madoTrayTarget;
    [menu addItem:madoTrayExitItem];

    madoTrayStatusItem.menu = menu;
    madoTrayStatusItem.length = NSSquareStatusItemLength;
    madoTrayStatusItem.visible = YES;
    madoTrayStatusItem.button.enabled = YES;
    madoTrayStatusItem.button.hidden = NO;
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
  });
}

void MadoTrayHide(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (madoTrayStatusItem == nil) {
      return;
    }

    [[NSStatusBar systemStatusBar] removeStatusItem:madoTrayStatusItem];
    [madoTrayStatusItem release];
    madoTrayStatusItem = nil;
    madoTrayShowItem = nil;
    madoTrayExitItem = nil;
  });
}

void MadoTraySetLocale(char* locale) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (locale == NULL) {
      madoTrayLocale = @"en";
    } else {
      NSString *value = [NSString stringWithUTF8String:locale];
      madoTrayLocale = [value isEqualToString:@"es"] ? @"es" : @"en";
    }
    MadoTrayUpdateMenuTexts();
  });
}
